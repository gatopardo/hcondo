package controller

import (
//	"log"
	"net/http"
        "fmt"
        "strings"
        "time"
        "strconv"

	"github.com/gatopardo/hcondo/app/model"
	"github.com/gatopardo/hcondo/app/shared/view"

        "github.com/gorilla/context"
	"github.com/josephspurrier/csrfbanana"
	"github.com/julienschmidt/httprouter"
  )

//------------------------------------------------
// EgrePerGET despliega formulario escoger periodo
func EgrePerGET(w http.ResponseWriter, r *http.Request) {
	sess     := model.Instance(r)
        params   := context.Get(r, "params").(httprouter.Params)
        sPerId   := params.ByName("pid")
	if sPerId == "0"{
	     lastPeriod, err    :=  model.PeriodTop(1)
	     putErr(err,sess, "Without last Period", "EgrePerGET top")
	      perId             :=  lastPeriod[0].Id
              sPerId             =  strconv.FormatInt(int64(perId), 10)
	}

        path := "/egreso/periodo/register/"+sPerId
        lisPeriod, err := model.Periods()

	putErr(err,sess, "No hay periodos", "EgrePerGET")
	v                  := view.New(r)
	v.Name              = "egreso/egresoper"
	v.Vars["token"]     = csrfbanana.Token(w, r, sess)
	v.Vars["Title"]     =  "Escoger Periodo"
        v.Vars["Action"]    =  path
        v.Vars["LisPeriod"] = lisPeriod
	v.Render(w)
 }
// ---------------------------------------------------
// EgrePerPOST procesa la forma enviada con periodo
func EgrePerPOST(w http.ResponseWriter, r *http.Request) {
        var egres model.EgresoN
        var period model.Periodo
        var  err  error
	sess     := model.Instance(r)
        sPerId   := r.FormValue("id")
	path     := "/egreso/register/"+sPerId
        action   := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
            var lisTipo []model.Tipo
            var lisEgre []model.EgresoN
            egres.PeriodId,  _   =  atoi32(r.FormValue("id"))
            period.Id            =  egres.PeriodId
            _                    =  (&period).PeriodById()
            egres.Period         =  period.Inicio
	    fecini               := egres.Period.Format(layout)
	    fecfin               := period.Final.Format(layout)
            lisTipo,  err        = model.Tipos("DB")
	    putErr(err,sess, "No hay tipos", "EgrePerGET")
            lisEgre, _          = (&egres).EgresPer()
	    v                  := view.New(r)
	    v.Name              = "egreso/egresoreg"
            v.Vars["token"]     = csrfbanana.Token(w, r, sess)
	    v.Vars["Title"]     =  "Crear Egreso"
            v.Vars["Action"]    = path
	    v.Vars["Inic"]      = fecini
	    v.Vars["Final"]     = fecfin
            v.Vars["Egreso"]    = egres
            v.Vars["LisTip"]    = lisTipo
            v.Vars["LisEgres"]  = lisEgre
            v.Render(w)
        }
        http.Redirect(w, r, "/", http.StatusFound)  //
 }
// ---------------------------------------------------
 func getFormEgre(e *  model.EgresoN, r *http.Request, b bool)(err error){
	   formato         :=  "2006/01/02"
           formato2        :=  "2006-01-02"
           e.TipoId, _      =  atoi32(r.FormValue("tipId"))
           e.PeriodId, _    =  atoi32(r.FormValue("periodId"))
           e.Period, _      =  time.Parse(formato,r.FormValue("period"))
	   if b{
               e.Fecha, _   =  time.Parse(formato2,r.FormValue("fecha"))
           }else{
               e.Fecha, _   =  time.Parse(formato,r.FormValue("fecha"))
            }
           e.Descripcion    =  r.FormValue("descripcion")
	   var unr int64
	   ramount         :=  r.FormValue("amount")
	   samount         :=   strings.ReplaceAll(ramount, ",","")
           unr, err         =  money2int64(samount)
           if err == nil {
                 e.Amount   =  unr
           }
       return
   }
// ---------------------------------------------------
// EgreRegPOST despliega formulario crear egreso
func EgreRegPOST(w http.ResponseWriter, r *http.Request) {
        var egres   model.EgresoN
        var period  model.Periodo
        var err  error
	sess     := model.Instance(r)
        sPerId   := r.FormValue("periodId")
	path     := "/egreso/list/"+sPerId
        action   := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
           getFormEgre(&egres, r, true)
           period.Id       =  egres.PeriodId
           err                 =  (&period).PeriodById()
           egres.PeriodId       =   period.Id
           err                 =  (&egres).EgresCreate()
           if err != nil { putErr(err,sess, "Guardando Egreso", "EgreRegPOST ")
//	       http.Redirect(w, r, "/", http.StatusFound)
               return
           } else {  // todo bien
                sess.AddFlash(view.Flash{"Egreso. creada: " , view.FlashSuccess})
           }
            var lisTipo []model.Tipo
            var lisEgre []model.EgresoN
            lisTipo, err  = model.Tipos("DB")
            if err != nil { putErr(err,sess, "No hay tipos", "EgreRegPOST ")}
            lisEgre,err           = (&egres).EgresPer()
            if err != nil { putErr(err,sess, "No hay egresos", "EgreRegPOST ")}
            v                  := view.New(r)
            v.Name              = "egreso/egresoreg"
            v.Vars["token"]     = csrfbanana.Token(w, r, sess)
            v.Vars["Title"]     = "Guardar Egreso"
            v.Vars["Action"]    = "/egreso/register/"+sPerId
            v.Vars["Egreso"]    = egres
            v.Vars["LisTip"]    = lisTipo
            v.Vars["LisEgres"]  = lisEgre
	    v.Render(w)
        }
       http.Redirect(w, r, path, http.StatusFound)
 }
// ---------------------------------------------------
// EgreUpGET despliega la pagina del usuario
func EgreUpGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
        var egres model.EgresoN
	var params httprouter.Params
	params         = context.Get(r, "params").(httprouter.Params)
	sId           := params.ByName("id")
	sPerId        :=  params.ByName("pid")
	id,_          := atoi32(sId)
        egres.Id       = id
        path          := "/egreso/update/"+ sId
	lisTipo,  err := model.Tipos("DB")
        if err != nil { putErr(err,sess, "No hay tipos", "EgreUpGET ") }
	err = (&egres).EgresById()
        if err != nil { putErr(err,sess, "No esta egreso", "EgreUpGET ") 
           path     = "/egreso/list/"+ sPerId
           sess.Save(r, w)
           http.Redirect(w, r, path, http.StatusFound)
           return
	}
	v                  := view.New(r)
	v.Name              = "egreso/egresoupdate"
	v.Vars["token"]     = csrfbanana.Token(w, r, sess)
        v.Vars["Title"]     =  "Actualizar Egreso"
        v.Vars["Action"]    =  path
        v.Vars["Egre"]      = egres
        v.Vars["LisTip"]    = lisTipo
        v.Render(w)
   }

// ---------------------------------------------------
 func   getEgreFormUp(e1, e2 model.EgresoN, r * http.Request)(stUp string){
        var sf string
        var sup []string
        formato        :=  "2006/01/02"
	if e1.TipoId != e2.TipoId {
             sf  =  fmt.Sprintf( " tipo_id = %d ", e1.TipoId )
	     sup = append(sup, sf)
	}

	if e1.Amount  != e2.Amount {
             sf  =  fmt.Sprintf( " amount = %d ", e1.Amount )
	     sup = append(sup, sf)
	}
        if e1.Fecha.Format(formato) != e2.Fecha.Format(formato) {
             sf  =  fmt.Sprintf( " fecha = '%s' ", e1.Fecha.Format(formato) )
	     sup = append(sup, sf)
	}
	if e1.Descripcion != e2.Descripcion {
             sf  =  fmt.Sprintf( " description = '%s' ", e1.Descripcion )
	     sup = append(sup, sf)
	}
        lon := len(sup)
        if lon  > 0 {
            sini :=  "update egresos set "
	    now         := time.Now()
	    sf           =  fmt.Sprintf( " ,  updated_at = '%s' ", now.Format(layout) )
            stUp  =  strings.Join(sup, ", ")
            sr   :=  fmt.Sprintf(" where egresos.id = %d ", e1.Id)
            stUp = sini + stUp + sf + sr
       }
         return
  }
// ---------------------------------------------------
// EgreUpPOST procesa la forma enviada con los datos
func EgreUpPOST(w http.ResponseWriter, r *http.Request) {
        var err error
        var eg,egres model.EgresoN
	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	sess        := model.Instance(r)
	sId         := params.ByName("id")
	sPerId      := r.FormValue("pid")
        Id,_        := atoi32(sId)
        egres.Id     = Id
        eg.Id        = Id
        action      := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
            err  = (&egres).EgresById()
            if err != nil { putErr(err,sess, "No esta egreso", "EgreUpPOST") }
	    getFormEgre(&eg,r, false)
	    st          :=  getEgreFormUp(eg, egres, r)
            if len(st) == 0{
                 sess.AddFlash(view.Flash{"No actualizacion solicitada", view.FlashSuccess})
            } else {
             err   =  egres.EgresUpdate(st)
             if err == nil{
                 sess.AddFlash(view.Flash{"Egreso actualizada exitosamente : " , view.FlashSuccess})
             } else  {  putErr(err,sess, "Error actualizando", "EgreUpPOST") }
		sess.Save(r, w)
           }
        }
        jmpToListE( sPerId, w,r  )
	return	
 }

//------------------------------------------------

//   jmpToListE( sPerId string, )
func   jmpToListE( sPerId string, w http.ResponseWriter,  r *http.Request) {
	var per model.Periodo
	sess            := model.Instance(r)
        path        := "/egreso/list/"+sPerId
        lisPeriod,err   := model.Periods()
        if err != nil { putErr(err,sess, "No esta egreso", "jmpToListE") 
            sess.Save(r, w)
         }
	 Id,_             := atoi32(sPerId)
           per.Id           = Id
           err              = (&per).PeriodById()
	   if err != nil { putErr(err,sess, "No Periodo", "jmpToListE") 
              sess.Save(r, w)
           }

        lisEgre, err         := model.EgresLim(Id)
        if err != nil { putErr(err,sess, "Listando Egreso", "jmpToListE") 
            sess.Save(r, w)
        }

        v                   := view.New(r)
	v.Name               = "egreso/egresolis"
	v.Vars["token"]      = csrfbanana.Token(w, r, sess)
	v.Vars["Action"]     = path
	v.Vars["Per"]        = per
        v.Vars["LisPeriod"]  = lisPeriod
        v.Vars["LisEgre"]    = lisEgre
        v.Vars["Level"]      =  sess.Values["level"]
	v.Render(w)
}

// ---------------------------------------------------
// EgreLisGET despliega formulario escoger periodo
func EgreLisGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
        var params httprouter.Params
        params      = context.Get(r, "params").(httprouter.Params)
	sPerId     :=  params.ByName("id")
	if sPerId == "0"{
	     lastPeriod, err :=  model.PeriodTop(1)
             if err != nil { putErr(err,sess, "No ultimo periodo", "EgreLisGET") }
	      perId             :=  lastPeriod[0].Id
              sPerId             =  strconv.FormatInt(int64(perId), 10)
	}

        lisPeriod, err := model.Periods()
	if err != nil { putErr(err,sess, "Obteniendo periodo", "CuotLisGET")
            sess.Save(r, w)
         }

	v                  := view.New(r)
	v.Name              = "egreso/egresoper"
	v.Vars["token"]     = csrfbanana.Token(w, r, sess)
        v.Vars["Title"]     =  "Listar"
	v.Vars["PerId"]     =  sPerId
        v.Vars["Action"]    =  "/egreso/list"+sPerId
        v.Vars["LisPeriod"] = lisPeriod
	v.Render(w)
 }
//------------------------------------------------
// EgreLis displays the egres page
func EgreLisPOST(w http.ResponseWriter, r *http.Request) {
        var Id  uint32
	var per  model.Periodo
	sess            := model.Instance(r)
	sPerId      := r.FormValue("pid")
        path        := "/egreso/list/"+sPerId
        action      := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
           lisPeriod,err    := model.Periods()
           if err != nil { putErr(err,sess, "Obteniendo Periodos", "EgreLisPOST") 
              sess.Save(r, w)
           }
           Id,_           = atoi32(r.FormValue("id"))
	   per.Id         = Id
 	   err            = (&per).PeriodById()
           if err != nil { putErr(err,sess, "Error con Periodos=", "EgreLisPOST") 
               sess.Save(r, w)
            }

           lisEgre, err   := model.EgresLim(Id)
           if err != nil { putErr(err,sess, "Error Listando Egresos=", "EgreLisPOST") 
              sess.Save(r, w)
            }
            v                   := view.New(r)
            v.Name               = "egreso/egresolis"
            v.Vars["token"]      = csrfbanana.Token(w, r, sess)
	    v.Vars["Action"]     = path
            v.Vars["Per"]        = per
            v.Vars["LisPeriod"]  = lisPeriod
            v.Vars["LisEgre"]    = lisEgre
            v.Vars["Level"]      =  sess.Values["level"]
            v.Render(w)
    }
      http.Redirect(w, r, "/cuota/list", http.StatusFound)
 }

//------------------------------------------------
// EgreDeleteGET handles the note deletion
 func EgreDeleteGET(w http.ResponseWriter, r *http.Request) {
        sess := model.Instance(r)
        var egres model.EgresoN
        var params httprouter.Params
        params       = context.Get(r, "params").(httprouter.Params)
	sPerId      := params.ByName("pid")
	sId         := params.ByName("id")
	id,_        := atoi32(sId)
        path        :=  "/egreso/delet3/"+sId
        egres.Id   = id
	err         := (&egres).EgresById()
	if err != nil { putErr(err,sess, "No hay egreso", "CuotDeleteGET")
	   path      = "/egreso/list/"+ sPerId
           sess.Save(r, w)
           http.Redirect(w, r, path, http.StatusFound)
           return
	}
	v                  := view.New(r)
	v.Name              = "egreso/egresodelete"
	v.Vars["token"]     = csrfbanana.Token(w, r, sess)
        v.Vars["Title"]     = "Eliminar Egreso"
        v.Vars["Action"]    = path
        v.Vars["Egre"]      = egres
	v.Render(w)
  }
// ---------------------------------------------------
// EgreDeletePOST procesa la forma enviada con los datos
func EgreDeletePOST(w http.ResponseWriter, r *http.Request) {
        var err error
        var egres model.Egreso
	sess := model.Instance(r)
	sPerId      := r.FormValue("pid")
	SId         := r.FormValue("id")
        Id,_        := atoi32(SId)
        egres.Id     = Id
        action      := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
             err = egres.EgresDelete()
	     if err != nil { putErr(err,sess, "Muy raro", "CuotDeletePOST")
              } else {
                  sess.AddFlash(view.Flash{"Egreso borrado!", view.FlashSuccess})
              }
              sess.Save(r, w)
        }
       jmpToListE( sPerId, w,r  )
 }
// ---------------------------------------------------
