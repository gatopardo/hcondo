package controller

import (
	"log"
	"net/http"
        "fmt"
        "strings"
        "time"
        "strconv"
       "encoding/json"

	"hcondo/app/model"
	"hcondo/app/shared/view"

        "github.com/gorilla/context"
	"github.com/josephspurrier/csrfbanana"
	"github.com/julienschmidt/httprouter"
  )
  var (
      formato         =  "2006/01/02"
      formato2        =  "2006-01-02"
  )
//------------------------------------------------
//------------------------------------------------
 func JIngreGET(w http.ResponseWriter, r *http.Request) {
	var peridi  model.Periodo
	var peridf  model.Periodo
	var lisIngre  []model.IngresoJ
	var  ingresoL  model.IngresoL
        var js []byte
        var params httprouter.Params
//	sess := model.Instance(r)
        params      = context.Get(r, "params").(httprouter.Params)
	sfec1       :=  params.ByName("fec1")[:7] + "-01"
	sfec2       :=  params.ByName("fec2")[:7] + "-01"
// fmt.Printf(" JIngreGET 0 %s %s\n", sfec1, sfec2)	
	dtfec1,err1  :=  time.Parse(layout, sfec1)
        if (err1 != nil  )  { log.Println("JIngretGET",err1) }
	dtfec2,err2  :=  time.Parse(layout, sfec2)
        if (err2 != nil  )  { log.Println("JIngretGET",err2) }
	if (err1 == nil) && (err2 == nil){
        dtfec1       =  time.Date(dtfec1.Year(), dtfec1.Month(),dtfec1.Day(), 0, 0, 0, 0, time.Local)
        err1         = (&peridi).PeriodByFec(dtfec1)
        if (err1 != nil  )  { log.Println("JIngretGET",err1) }
        dtfec2       =  time.Date(dtfec2.Year(), dtfec2.Month(),dtfec2.Day(), 0, 0, 0, 0, time.Local)
        err2         = (&peridf).PeriodByFec(dtfec2)
        if (err2 != nil  )  { log.Println("JIngretGET",err2) }
	if (err1 == nil) && (err2 == nil){
          lisIngre, err1  = model.IngresoJPer( peridi.Inicio, peridf.Final )
	  var  num int64 = 0
// fmt.Printf(" JIngreGET 1 %d \n", len(lisIngre))	
          if err1 != nil { log.Println(err1) 
          }else {
            if len(lisIngre) == 0 {
	        var dfec time.Time
                dfec,_ = time.Parse(formato , "1900/01/01")
		ingresoj := model.IngresoJ{Tipo: "X"  , Fecha: dfec, Amount: num, Descripcion: ""}
                lisIngre =  append( lisIngre, ingresoj)
            }
            ingresoL.Period  =  peridi.Inicio
            ingresoL.LisIngre =  lisIngre
            js, err1 =  json.Marshal(ingresoL)
            if err1 == nil{
               w.Header().Set("Content-Type", "application/json")
               w.Write(js)
	       return
            }
           }
          }
          }
          log.Println("JIngreGET  ", err1)
          http.Error(w, err1.Error(), http.StatusInternalServerError)
          return
 }



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











