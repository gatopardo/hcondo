package controller

import (
//	"log"
	"net/http"
        "fmt"
        "strings"
        "strconv"
        "time"

	"github.com/gatopardo/hcondo/app/model"
	"github.com/gatopardo/hcondo/app/shared/view"

        "github.com/gorilla/context"
	"github.com/josephspurrier/csrfbanana"
	"github.com/julienschmidt/httprouter"
  )
  // -------------------------------------------------:--
// CuotPerGET despliega formulario escoger periodo
func CuotPerGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
	params   := context.Get(r, "params").(httprouter.Params)
        sPerId   := params.ByName("pid")
	if sPerId == "0"{
	     lastPeriod, err :=  model.PeriodTop(1)
	     putErr(err,sess, "Without last Period", "CuotPerGET top")
	      perId             :=  lastPeriod[0].Id
              sPerId             =  strconv.FormatInt(int64(perId), 10)
	}
        path := "/cuota/periodo/register/"+sPerId

        lisPeriod, err := model.Periods()
        putErr(err,sess, "No hay periodos", "CuotPerGET" )
	v                  := view.New(r)
	v.Name              = "cuota/cuotper"
	v.Vars["token"]     = csrfbanana.Token(w, r, sess)
        v.Vars["Title"]     = "Escoger Periodo"
        v.Vars["Action"]    = path
        v.Vars["LisPeriod"] = lisPeriod
	v.Render(w)
 }
 //----------------------------------------------------
// CuotPerPOST procesa la forma enviada con periodo
func CuotPerPOST(w http.ResponseWriter, r *http.Request) {
        var cuot model.CuotaN
        var period model.Periodo
	sess          := model.Instance(r)
        sPerId        := r.FormValue("id")
	path          := "/cuota/register/"+sPerId
        action        := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
            var lisTipo []model.Tipo
            var lisCuot []model.CuotaN
            pid,  _            :=  atoi32(sPerId)
	    period.Id           =  pid
	    cuot.PeriodId       =  pid
            cuot.Transfe        =  true
            _                   =  (&period).PeriodById()
            cuot.Period         =  period.Inicio
	    fecini             := cuot.Period.Format(layout)
	    fecfin             := period.Final.Format(layout)
            lisApts, err       :=  model.Apts()
	    putErr(err,sess, "No hay aptos", "CuotPerGET ")
            lisTipo,  err       = model.Tipos("CR")
	     putErr(err,sess, "No hay tipos", "CuotPerGET ")
             lisCuot, _         = (&cuot).CuotsPer()
	     v                 := view.New(r)
	     v.Name             = "cuota/cuotreg"
             v.Vars["token"]    = csrfbanana.Token(w, r, sess)
	     v.Vars["Title"]    = "Crear Cuota"
             v.Vars["Action"]   = path
             v.Vars["Cuot"]     = cuot
	     v.Vars["Inic"]     = fecini
	     v.Vars["Final"]    = fecfin
             v.Vars["LisApt"]   = lisApts
             v.Vars["LisTip"]   = lisTipo
             v.Vars["LisCuots"] = lisCuot
             v.Render(w)
        }
       http.Redirect(w, r, "/", http.StatusFound)
 }
// ---------------------------------------------------
 func getFormCuot(c *  model.CuotaN, r *http.Request, b bool)(err error){
//           var t bool
           formato         :=  "2006/01/02"
           formato2        :=  "2006-01-02"
           c.ApartaId, _    =  atoi32(r.FormValue("aptId"))
           c.TipoId, _      =  atoi32(r.FormValue("tipId"))
           c.PeriodId, _    =  atoi32(r.FormValue("periodId"))
           stPeriod        := r.FormValue("period")
           stFecha         := r.FormValue("fecha")
           c.Period,_       =  time.Parse(formato,stPeriod)
	   if b {
              c.Fecha, _    =  time.Parse(formato2,stFecha)
            }else{
              c.Fecha, _    =  time.Parse(formato,stFecha)
	    }
           c.Transfe,_      = strconv.ParseBool(r.FormValue("transfe"))
           ramount         :=  r.FormValue("amount")
           samount         :=   strings.ReplaceAll(ramount, ",","")
           unr, err        :=  money2int64(samount)
           if err == nil {
                 c.Amount   =  unr
            }
       return
   }
// ---------------------------------------------------
// CuotRegPOST despliega formulario crear cuota
func CuotRegPOST(w http.ResponseWriter, r *http.Request) {
        var cuot   model.CuotaN
        var period  model.Periodo
        var err  error
        sPerId   := r.FormValue("periodId")
	path     := "/cuota/list/"+sPerId
	sess   := model.Instance(r)
        action        := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
           getFormCuot(&cuot, r, true)
           period.Id           =  cuot.PeriodId
           err                 =  (&period).PeriodById()
	   putErr(err,sess, "No hay periodos", "CuotRegPOST ")
           err                 =  (&cuot).CuotCreate()
	   if err != nil { putErr(err,sess, "No hay create", "CuotRegPOST ")
               return
           } else {  sess.AddFlash(view.Flash{"Cuota creada", view.FlashSuccess}) }
            var lisApto []model.Aparta
            var lisTipo []model.Tipo
            var lisCuot []model.CuotaN
            lisApto, err  = model.Apts()
	    if err != nil { putErr(err,sess, "No hay apartas", "CuotRegPOST ")}
            lisTipo, err  = model.Tipos("CR")
	    if err != nil { putErr(err,sess, "No hay tipos", "CuotRegPOST ")}
            lisCuot, err         = (&cuot).CuotsPer()
	    if err != nil { putErr(err,sess, "No hay cuotas", "CuotRegPOST ")}
            v                  := view.New(r)
            v.Name              = "cuota/cuotreg"
            v.Vars["token"]     = csrfbanana.Token(w, r, sess)
	    v.Vars["Title"]     =  "Guardar Cuota"
            v.Vars["Action"]    =  "/cuota/register/"+sPerId
            v.Vars["Cuot"]      = cuot
            v.Vars["LisApt"]    = lisApto
            v.Vars["LisTip"]    = lisTipo
            v.Vars["LisCuots"]  = lisCuot
	    v.Render(w)
        }
	http.Redirect(w, r, path, http.StatusFound)
 }
// ---------------------------------------------------
// CuotUpGET despliega la pagina del usuario
func CuotUpGET(w http.ResponseWriter, r *http.Request) {
        var lisTipo []model.Tipo
	sess := model.Instance(r)
        var cuot model.CuotaN
	var params httprouter.Params
	params    = context.Get(r, "params").(httprouter.Params)
	sId      :=  params.ByName("id")
	sPerId   :=  params.ByName("pid")
	id,_     := atoi32(sId)
        cuot.Id   = id
        path     :=  "/cuota/update/" + sId
        lisApts, err    :=  model.Apts()
        if err != nil { putErr(err,sess, "No hay aptos", "CuotUpGET ")}
        lisTipo,  err    = model.Tipos("CR")
        if err != nil { putErr(err,sess, "No hay tipos", "CuotUpGET ")}
	err = (&cuot).CuotById()
        if err != nil { putErr(err,sess, "No hay cuotas", "CuotUpGET ")
	   path     = "/cuota/list/" + sPerId
           sess.Save(r, w)
           http.Redirect(w, r, path, http.StatusFound)
           return
	}
	v                    :=  view.New(r)
	v.Name                =  "cuota/cuotupdate"
	v.Vars["token"]       =  csrfbanana.Token(w, r, sess)
	v.Vars["Title"]       =  "Actualizar Cuota"
        v.Vars["Action"]      =  path
        v.Vars["Cuot"]        =  cuot
        v.Vars["LisApt"]      =  lisApts
        v.Vars["LisTip"]      =  lisTipo
        v.Render(w)
   }

// ---------------------------------------------------
 func   getCuotFormUp(c1,c2 model.CuotaN, r * http.Request)(stUp string){
        var sf string
	var sup  []string
        formato        :=  "2006/01/02"
	if c1.ApartaId != c2.ApartaId {
             sf  =  fmt.Sprintf( " aparta_id = %d ", c1.ApartaId )
	     sup = append(sup, sf)
	}

	if c1.TipoId != c2.TipoId {
             sf  =  fmt.Sprintf( " tipo_id = %d ", c1.TipoId )
	     sup = append(sup, sf)
	}

	if c1.Fecha.Format(formato) != c2.Fecha.Format(formato) {
             sf  =  fmt.Sprintf( " fecha = '%s' ", c1.Fecha.Format(formato) )
	     sup = append(sup, sf)
	}

	if c1.Amount != c2.Amount {
             sf  =  fmt.Sprintf( " amount = %d ", c1.Amount )
	     sup = append(sup, sf)
	}
	if c1.Transfe != c2.Transfe {
              sf =  fmt.Sprintf(" transfe = %t ", c1.Transfe )
	     sup = append(sup, sf)
	}
       lon := len(sup)
       if lon  > 0 {
	    now         := time.Now()
	    sf           =  fmt.Sprintf( " , updated_at = '%s' ", now.Format(formato) )
            sini        :=  "update cuotas set "
            stUp         =  strings.Join(sup, ", ")
	    sr          :=  fmt.Sprintf(" where cuotas.id = %d ", c1.Id)

            stUp         = sini + stUp + sf +  sr
       }
         return
  }
  // ---------------------------------------------------
// CuotUpPOST procesa la forma enviada con los datos
func CuotUpPOST(w http.ResponseWriter, r *http.Request) {
        var err error
        var c, cuot model.CuotaN
        var params httprouter.Params
	params       = context.Get(r, "params").(httprouter.Params)
	sess        := model.Instance(r)
	SId         := params.ByName("id")
        sPerId      := r.FormValue("pid")

	Id,_        := atoi32(SId)
        cuot.Id      = Id
        c.Id         = Id
        action      := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
	    err = (&cuot).CuotById()
            if err != nil { putErr(err,sess, "No hay cuota", "CuotUpPOST ")}
	    getFormCuot(&c,r, false)
            st          :=  getCuotFormUp(c, cuot, r)
            if len(st) == 0{
                 sess.AddFlash(view.Flash{"No actualizacion solicitada", view.FlashSuccess})
            } else {
             err   =  cuot.CuotUpdate(st)
             if err == nil{
                 sess.AddFlash(view.Flash{"Cuota actualizada exitosamente : " , view.FlashSuccess})
             } else   { putErr(err,sess, "Errror actualizando", "CuotUpPOST ")}
		sess.Save(r, w)
           }
        }
        jmpToList( sPerId, w,r  )
	return
 }
//------------------------------------------------
//   jmpToList( sPerId string, )
func   jmpToList( sPerId string, w http.ResponseWriter,  r *http.Request) {
	var per model.Periodo
	sess            := model.Instance(r)
        path            := "/cuota/list/"+sPerId
        lisPeriod,err   := model.Periods()
        if err != nil { putErr(err,sess, "No hay periodo", "jmpToList ")
            sess.Save(r, w)
         }
	 Id,_              := atoi32(sPerId)
           per.Id           = Id
           err              = (&per).PeriodById()
           if err != nil { putErr(err,sess, "No hay periodo", "jmpToList ")
               sess.Save(r, w)
          }

        lisCuot, err       := model.CuotLim(Id)
        if err != nil { putErr(err,sess, "istando cuotas ", "jmpToList ")
            sess.Save(r, w)
        }

        v                   := view.New(r)
	v.Name               = "cuota/cuotlis"
	v.Vars["token"]      = csrfbanana.Token(w, r, sess)
	v.Vars["Per"]        = per
	v.Vars["PerId"]      = sPerId
	v.Vars["Action"]     = path
        v.Vars["LisPeriod"]  = lisPeriod
        v.Vars["LisCuot"]    = lisCuot
        v.Vars["Level"]      = sess.Values["level"]
	v.Render(w)
}


//------------------------------------------------
// CuotLisGET displays the cuot page
func CuotLisGET(w http.ResponseWriter, r *http.Request) {
	sess            := model.Instance(r)
	var params httprouter.Params
        params      = context.Get(r, "params").(httprouter.Params)
	sPerId     :=  params.ByName("id")
	if sPerId == "0"{
	     lastPeriod, err :=  model.PeriodTop(1)
             if err != nil { putErr(err,sess, "No ultimo periodo", "CuotLisGET") }
	      perId             :=  lastPeriod[0].Id
              sPerId             =  strconv.FormatInt(int64(perId), 10)
	}
        lisPeriod,err   := model.Periods()
        if err != nil { putErr(err,sess, "Obteniendo periodo", "CuotLisGET")
            sess.Save(r, w)
         }

         v                   := view.New(r)
         v.Name               = "cuota/cuotper"
         v.Vars["token"]      = csrfbanana.Token(w, r, sess)
         v.Vars["Title"]      =  "Listar"
         v.Vars["PerId"]      =  sPerId
         v.Vars["Action"]     =  "/cuota/list/"+sPerId
         v.Vars["LisPeriod"]  = lisPeriod
         v.Render(w)
 }

//------------------------------------------------
// CuotLis displays the cuot page
func CuotLisPOST(w http.ResponseWriter, r *http.Request) {
        var Id  uint32
	var per model.Periodo
	sess            := model.Instance(r)
	sPerId      := r.FormValue("pid")
        path        := "/cuota/list/"+sPerId
        action      := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
           lisPeriod,err    := model.Periods()
           if err != nil { putErr(err,sess, "No hay periodos", "CuotLisPOST")
               sess.Save(r, w)
            }
           sId             := r.FormValue("id")
           Id,_             = atoi32(sId)
           per.Id           = Id
           err              = (&per).PeriodById()
           if err != nil { putErr(err,sess, "No hay periodo", "CuotLisPOST")
            sess.Save(r, w)
           }

        lisCuot, err         := model.CuotLim(Id)
        if err != nil { putErr(err,sess, "Listando cuotas", "CuotLisPOST")
            sess.Save(r, w)
        }
	v                   := view.New(r)
	v.Name               = "cuota/cuotlis"
	v.Vars["token"]      = csrfbanana.Token(w, r, sess)
	v.Vars["Per"]        = per
	v.Vars["PerId"]      = sPerId
	v.Vars["Action"]     = path
        v.Vars["LisPeriod"]  = lisPeriod
        v.Vars["LisCuot"]    = lisCuot
        v.Vars["Level"]      = sess.Values["level"]
	v.Render(w)
       }
       http.Redirect(w, r, path, http.StatusFound)
 }

//------------------------------------------------
// CuotDeleteGET handles the note deletion
 func CuotDeleteGET(w http.ResponseWriter, r *http.Request) {
        sess := model.Instance(r)
        var cuot model.CuotaN
        var params httprouter.Params
        params       = context.Get(r, "params").(httprouter.Params)
	sPerId      := params.ByName("pid")
	sId         := params.ByName("id")
	id,_        := atoi32(sId)
        path        := "/cuota/delete/"+sId
        cuot.Id      = id
	err         := (&cuot).CuotById()
        if err != nil { putErr(err,sess, "No hay cuota", "CuotDeleteGET")
	   path      = "/cuota/list/"+ sPerId
           sess.Save(r, w)
           http.Redirect(w, r, path, http.StatusFound)
           return
	}
	v                    := view.New(r)
	v.Name                = "cuota/cuotdelete"
	v.Vars["token"]       = csrfbanana.Token(w, r, sess)
        v.Vars["Title"]       = "Eliminar Cuota"
        v.Vars["Action"]      = path
        v.Vars["Cuot"]        = cuot
	v.Render(w)
	return
  }

// ---------------------------------------------------
// CuotDeletePOST procesa la forma enviada con los datos
func CuotDeletePOST(w http.ResponseWriter, r *http.Request) {
        var err error
        var cuot model.Cuota
	sess := model.Instance(r)
	sPerId      := r.FormValue("pid")
	SId         := r.FormValue("id")
        Id,_        := atoi32(SId)
        cuot.Id      = Id
        action      := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
             err = cuot.CuotDelete()
             if err != nil { putErr(err,sess, "No borrada cuota", "CuotDeletePOST")
              } else {
                  sess.AddFlash(view.Flash{"Cuota borrado!", view.FlashSuccess})
              }
              sess.Save(r, w)
        }
        jmpToList( sPerId, w,r  )
	return
 }



