package controller

import (
	"log"
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
  var (
      formato         =  "2006/01/02"
      formato2        =  "2006-01-02"
  )
//------------------------------------------------
// ---------------------------------------------------
// IngrePerGET despliega formulario escoger periodo
func IngrePerGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
	params   := context.Get(r, "params").(httprouter.Params)
        sPerId   := params.ByName("pid")
	if sPerId == "0"{
              lastPeriod, err :=  model.PeriodTop(1)
	     putErr(err,sess, "Without last Period", "IngrePerGET")
	      perId             :=  lastPeriod[0].Id
              sPerId             =  strconv.FormatInt(int64(perId), 10)
	}
        path := "/ingreso/periodo/register/"+sPerId

        lisPeriod, err := model.Periods()
        putErr(err,sess, "No hay periodos", "IngrePerGET ")
	v                  := view.New(r)
	v.Name              = "ingreso/ingresoper"
	v.Vars["token"]     = csrfbanana.Token(w, r, sess)
        v.Vars["Title"]     =  "Ingreso"
        v.Vars["Action"]    =  path
        v.Vars["LisPeriod"] = lisPeriod
	v.Render(w)
 }
// ---------------------------------------------------
// IngrePerPOST procesa la forma enviada con periodo
func IngrePerPOST(w http.ResponseWriter, r *http.Request) {
        var ingres model.IngresoN
        var period model.Periodo
        var  err  error
	sess          := model.Instance(r)
	sPerId        := r.FormValue("id")
	path          := "/ingreso/register/" + sPerId
        action        := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
            pid,_                :=  atoi32(sPerId)
            ingres.PeriodId       =  pid
            period.Id             =  pid
            _                     =  (&period).PeriodById()
            ingres.Period         =  period.Inicio
	    fecini               :=  ingres.Period.Format(layout)
	    fecfin               :=  period.Final.Format(layout)
            var lisTipo []model.Tipo
            var lisIngre []model.IngresoN
            lisTipo,  err        = model.Tipos("CR")
	    putErr(err,sess, "No hay tipos", "IngrePerGET ")
            lisIngre, _          = model.IngresLim(ingres.PeriodId)
	    v                   := view.New(r)
	    v.Name               = "ingreso/ingresoreg"
            v.Vars["token"]      = csrfbanana.Token(w, r, sess)
	    v.Vars["Title"]      = "Crear Ingreso"
            v.Vars["Action"]     = path
	    v.Vars["Inic"]       = fecini
	    v.Vars["Final"]      = fecfin
            v.Vars["Ingreso"]    = ingres
            v.Vars["LisTip"]     = lisTipo
            v.Vars["LisIngres"]  = lisIngre
            v.Render(w)
        }
	http.Redirect(w, r, "/", http.StatusFound)
 }
// ---------------------------------------------------
 func getFormIngre(ing *  model.IngresoN, r *http.Request, b bool)(err error){
           formato         :=  "2006/01/02"
           formato2        :=  "2006-01-02"
	   ing.TipoId, _    =  atoi32(r.FormValue("tipId"))
           ing.PeriodId, _  =  atoi32(r.FormValue("periodId"))
           ing.Period, _    =  time.Parse(layout,r.FormValue("period"))
	   if b{
               ing.Fecha, _ =  time.Parse(formato2,r.FormValue("fecha"))
           }else{
               ing.Fecha, _ =  time.Parse(formato,r.FormValue("fecha"))
           }
	   ing.Descripcion   =  r.FormValue("descripcion")
	   var nro int64
	   ramount         :=  r.FormValue("amount")
	   samount         :=   strings.ReplaceAll(ramount, ",","")
           nro, err       = money2int64(samount)
           if err == nil {
		   ing.Amount   =  nro
            }
           ing.Descripcion   =  r.FormValue("descripcion")
       return
   }
// ---------------------------------------------------
// IngreRegPOST despliega formulario crear ingreso
func IngreRegPOST(w http.ResponseWriter, r *http.Request) {
        var ingres   model.IngresoN
        var period  model.Periodo
        var err  error
	sess          := model.Instance(r)
	sPerId        := r.FormValue("periodId")  
	path          := "/ingreso/list/"+sPerId
        action        := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
           getFormIngre(&ingres, r, true)
           period.Id            =  ingres.PeriodId
           _                    =  (&period).PeriodById()
           err                  =  (&ingres).IngresCreate()
           if err != nil{  putErr(err,sess, "No hay tipos", "IngrePerPOST ")
               return
//	       http.Redirect(w, r, "/ingreso/list", http.StatusFound)
            } else {  // todo bien
                sess.AddFlash(view.Flash{"Ingreso. creado: " , view.FlashSuccess}) 
            }

            var lisTipo []model.Tipo
            var lisIngre []model.IngresoN
            lisTipo, err  = model.Tipos("CR")
            putErr(err,sess, "No hay tipos", "IngrePerGET ")
            lisIngre, _           = model.IngresLim(period.Id)
            v                    := view.New(r)
            v.Name                = "ingreso/ingresoreg"
            v.Vars["token"]       = csrfbanana.Token(w, r, sess)
            v.Vars["Ingreso"]     = ingres
            v.Vars["LisTip"]      = lisTipo
            v.Vars["LisIngres"]   = lisIngre
	    v.Render(w)
        }
	http.Redirect(w, r, path, http.StatusFound)
 }
// ---------------------------------------------------
// IngreUpGET despliega la pagina del usuario
func IngreUpGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
        var ingres model.IngresoN
	var params httprouter.Params
	params  = context.Get(r, "params").(httprouter.Params)
	id,_   := atoi32(params.ByName("id"))
        path   :=  "/ingreso/list"
        ingres.Id = id
	lisTipo,  err        := model.Tipos("CR")
            if err != nil {
                 sess.AddFlash(view.Flash{"No hay tipos ", view.FlashError})
            }
	err = (&ingres).IngresById()
	if err != nil { // Si no existe el ingreso
           log.Println(err)
           sess.AddFlash(view.Flash{"Es raro. No esta ingreso.", view.FlashError})
           sess.Save(r, w)
           http.Redirect(w, r, path, http.StatusFound)
           return
	}
	v                    := view.New(r)
	v.Name                = "ingreso/ingresoupdate"
	v.Vars["token"]       = csrfbanana.Token(w, r, sess)
        v.Vars["Title"]     =  "Actualizar Ingreso"
        v.Vars["Action"]    =  "/ingreso/update"
        v.Vars["Ingre"]       = ingres
        v.Vars["LisTip"]      = lisTipo
        v.Render(w)
   }

// ---------------------------------------------------
 func   getIngreFormUp(i1, i2 model.IngresoN ,r * http.Request)(stUp string){
        var sf string
        var sup []string
        formato        :=  "2006/01/02"

	if i1.TipoId != i2.TipoId {
             sf  =  fmt.Sprintf( " tipo_id = %d ", i1.TipoId )
	     sup = append(sup, sf)
	}
	if i1.Amount  != i2.Amount {
             sf  =  fmt.Sprintf( " amount = '%d' ", i1.Amount )
	     sup = append(sup, sf)
	}
        if i1.Fecha.Format(formato) != i2.Fecha.Format(formato) {
             sf  =  fmt.Sprintf( " fecha = '%s' ", i1.Fecha.Format(layout) )
	     sup = append(sup, sf)
	}
	if i1.Descripcion != i2.Descripcion {
             sf  =  fmt.Sprintf( " descripcion = '%s' ", i1.Descripcion )
	     sup = append(sup, sf)
	}
        lon := len(sup)
        if lon  > 0 {
            sini :=  "update ingresos set "
	    now         := time.Now()
	    sf           =  fmt.Sprintf( " , updated_at = '%s' ", now.Format(layout) )
            stUp  =  strings.Join(sup, ", ")
            sr   :=  fmt.Sprintf(" where ingresos.id = %d ", i1.Id)
            stUp = sini + stUp + sf + sr
       }
       return
  }
// ---------------------------------------------------
// IngreUpPOST procesa la forma enviada con los datos
func IngreUpPOST(w http.ResponseWriter, r *http.Request) {
        var err error
        var ing, ingres model.IngresoN
	sess := model.Instance(r)
        var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	SId         := params.ByName("id")
        Id,_        := atoi32(SId)
        ingres.Id    = Id
        ing.Id       = Id
        path        :=  "/ingreso/list"
        action      := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
            err  = (&ingres).IngresById()
	    if err != nil { // Si no existe ingreso
                  sess.AddFlash(view.Flash{"Es raro. No esta ingreso.", view.FlashError})
            }
	    getFormIngre(&ing,r,false)
            st          :=  getIngreFormUp(ing, ingres, r)

            if len(st) == 0{
                 sess.AddFlash(view.Flash{"No actualizacion solicitada", view.FlashSuccess})
            } else {
             err   =  ingres.IngresUpdate(st)
             if err == nil{
                 sess.AddFlash(view.Flash{"Ingreso actualizada exitosamente : " , view.FlashSuccess})
             } else       {
		log.Println(err)
		sess.AddFlash(view.Flash{"Un error ocurrio actualizando.", view.FlashError})
	     }
		sess.Save(r, w)
           }
        }
	http.Redirect(w, r, path, http.StatusFound)
 }
// ---------------------------------------------------
// IngreLisGET despliega formulario escoger periodo
func IngreLisGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
        lisPeriod, err := model.Periods()
        if err != nil {
             sess.AddFlash(view.Flash{"No hay periodos ", view.FlashError})
         }
	v                  := view.New(r)
	v.Name              = "ingreso/ingresoper"
	v.Vars["token"]     = csrfbanana.Token(w, r, sess)
        v.Vars["LisPeriod"] = lisPeriod
        v.Vars["Title"]     =  "Listar"
        v.Vars["Action"]    =  "/ingreso/list"
	v.Render(w)
 }
//------------------------------------------------
// IngreLis displays the ingres page
func IngreLisPOST(w http.ResponseWriter, r *http.Request) {
        var Id  uint32
	var per model.Periodo
	sess            := model.Instance(r)
        lisPeriod,err    := model.Periods()
        if err != nil {
            log.Println(err)
	    sess.AddFlash(view.Flash{"Error con Lista Periodos.", view.FlashError})
            sess.Save(r, w)
         }
         Id,_             = atoi32(r.FormValue("id"))
	per.Id                 = Id
	err = (&per).PeriodById();
        if err != nil {
            log.Println(err)
	    sess.AddFlash(view.Flash{"Error con  Periodo.", view.FlashError})
            sess.Save(r, w)
         }
        lisIngre, err         := model.IngresLim(Id)
        if err != nil {
            log.Println(err)
	    sess.AddFlash(view.Flash{"Error Listando Ingresos.", view.FlashError})
            sess.Save(r, w)
         }
	v                   := view.New(r)
	v.Name               = "ingreso/ingresolis"
	v.Vars["token"]      = csrfbanana.Token(w, r, sess)
	v.Vars["Per"]        = per
        v.Vars["LisPeriod"]  = lisPeriod
        v.Vars["LisIngre"]    = lisIngre
        v.Vars["Level"]      =  sess.Values["level"]
	v.Render(w)
 }

//------------------------------------------------
// IngreDeleteGET handles the note deletion
 func IngreDeleteGET(w http.ResponseWriter, r *http.Request) {
        sess := model.Instance(r)
        var ingres model.IngresoN
        var params httprouter.Params
        params = context.Get(r, "params").(httprouter.Params)
	SId         := params.ByName("id")
	id,_        := atoi32(SId)
        path        :=  "/ingreso/list"
        ingres.Id   = id
	err         := (&ingres).IngresById()
	if err != nil { // Si no existe ingreso
           log.Println(err)
           sess.AddFlash(view.Flash{"Es raro. No hay ingreso.", view.FlashError})
           sess.Save(r, w)
           http.Redirect(w, r, path, http.StatusFound)
           return
	}
	v                    := view.New(r)
	v.Name                = "ingreso/ingresodelete"
	v.Vars["token"]       = csrfbanana.Token(w, r, sess)
        v.Vars["Title"]     =  "Eliminar Ingreso"
        v.Vars["Action"]    =  "/ingreso/delete"
        v.Vars["Ingre"]        = ingres
        v.Vars["Level"]       =  sess.Values["level"]
	v.Render(w)
  }

// ---------------------------------------------------
// IngreUpPOST procesa la forma enviada con los datos
func IngreDeletePOST(w http.ResponseWriter, r *http.Request) {
        var err error
        var ingres model.Ingreso
	sess := model.Instance(r)
        var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	SId         := params.ByName("id")
        Id,_        := atoi32(SId)
        ingres.Id    = Id
        path        :=  "/ingreso/list"
        action      := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
             err = ingres.IngresDelete()
             if err != nil {
                 log.Println(err)
                 sess.AddFlash(view.Flash{"Error no posible. Auxilio.", view.FlashError})
              } else {
                  sess.AddFlash(view.Flash{"Ingreso borrado!", view.FlashSuccess})
              }
              sess.Save(r, w)
        }
	http.Redirect(w, r, path, http.StatusFound)
 }
// ---------------------------------------------------











