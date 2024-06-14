package controller

import (
        "log"
	"net/http"
        "strings"
        "fmt"
        "time"
        "encoding/json"

	"hcondo/app/model"
	"hcondo/app/shared/view"

        "github.com/gorilla/context"
	"github.com/josephspurrier/csrfbanana"
	"github.com/julienschmidt/httprouter"
  )
//      Refill any form fields
// view.Repopulate([]string{"name"}, r.Form, v.Vars)

// ---------------------------------------------------
  func procesaApt(lsAptTots []model.AmtAptTot)(pos int, scuot, smonto int64){
         var  sdif int64
	 scuot  = int64(0)
	 smonto = int64(0)
	 sdif   =  int64(0)
	 sfec  := time.Date(1,1,1,0,0,0,0,time.UTC)
         pos = 0
	 j := 0
	 k := 0
         for i, amt := range lsAptTots {
             k = i
	     smonto = smonto +  amt.Monto
	     if sfec != amt.Inicio {
	          scuot = scuot + amt.Cuota
	     }
	     sdif =  scuot - smonto
	     if sdif <= int64(0){
		     j = pos
		     pos = i
	     }
	     lsAptTots[i].Dif = sdif
	     sfec = amt.Inicio
	 }
         if  (k - pos) <  5{
             pos = j
	     if  (j  - pos) < 5 && (k - 10) > 0{
                pos = k - 10
	     }
	 }
        return
 }
// ---------------------------------------------------
// japtget json service for apt state
func JAptGET(w http.ResponseWriter, r *http.Request) {
	var peridf  model.Periodo
        var params httprouter.Params
        var jpers  model.Jperson
	var aparta model.Aparta
	var aptEstadL AptEstadL
	var aptEstad []AptEstadJ
        var aptDet AptEstadJ
	
	//        var lisPaym []model.CuotApt
//	var arPaym ArPay
//	var dt11, dt22  time.Time
	var dt22  time.Time
	var err  error
	sess := model.Instance(r)
        params           = context.Get(r, "params").(httprouter.Params)
//	sfec1           :=  params.ByName("fec1")[:7]+"-01"
        sfec1           := "0001-01-01"
	sfec2           :=  params.ByName("fec1")[:7]+"-01"
	dIni,_          := time.Parse( layout, sfec1)
//	dFini,_         := time.Parse(layout, sfec2)
        sId             :=  params.ByName("id")
        uid,_           :=  atoi32(sId)
 fmt.Printf("JAptGET 0 %s %s\n", sfec1, sfec2)
//	dt11,err         =  time.Parse(layout, sfec1)
//        if err == nil {
	    dt22,err   =  time.Parse(layout, sfec2)
//            err        =  (&peridi).PeriodByFec(dt11)
//	    if err    ==  nil {
               err     =  (&peridf).PeriodByFec(dt22)
//            }
//        }
        if err      != nil {
	        log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	        return
        }	
	_, err         =   (&jpers).JPersByUserId(uid)
	if err == model.ErrNoResult {
           log.Println("JAPTGET ", err)
	        http.Error(w, err.Error(), http.StatusBadRequest)
	   return
        }
	aparta.Id  = jpers.AptId
        err = (&aparta).AptById()
        if err      != nil {
	        log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	        return
        }

        fmt.Printf("JAptGET 1 %s %s\n", peridf.Final.Format(layout), dIni.Format(layout))
/*        lisPaym, err   =  model.Payments(jpers.AptId, peridf.Inicio, peridi.Inicio)
        if err == nil {
           arPaym.Apto   = jpers.Apto
	   arPaym.Final  = peridf.Final
	   arPaym.APaym  = lisPaym
*/
        amtTots, err := model.AptDetails(aparta.Codigo, dIni , peridf.Final )

// fmt.Printf( "JAptGET 2 %d %s %s %s\n", len(amtTots), aparta.Codigo,dIni.Format(layout), peridf.Final.Format(layout))
 // fmt.Println( amtTots)
  	if err != nil {
		log.Println(err)
		return
	}

        p,scuot,smonto := procesaApt(amtTots)
	aptEstadL.Apt     =  aparta.Codigo
	aptEstadL.Period  =  peridf.Final
	aptEstadL.SCuota  =  scuot
	aptEstadL.SAmount =  smonto
	lon :=  len(amtTots)
	for i  := lon - 1; i >= p; i-- {
		aptDet = AptEstadJ{ Fecha:   amtTots[i].Fecha,
                                    Cuota:   amtTots[i].Cuota,
                                    Amount:  amtTots[i].Monto,
                                    Balance: amtTots[i].Dif,
	                          }
	   aptEstad  = append(aptEstad, aptDet)
         }
	 aptEstadL.LisEstad = aptEstad
// fmt.Printf( "JAptGET 3 %d %s %s %d %d\n", len(aptEstad), aparta.Codigo,peridf.Final.Format(layout), smonto/100, scuot/100)
 fmt.Println( "JAptGET 4 \n", aptEstad[:5] )

           var js []byte
           js, err =  json.Marshal(aptEstadL)
// fmt.Println( "JAptGET 5 \n", js )
           if err == nil{
               w.Header().Set("Content-Type", "application/json")
               w.Write(js)
	       return
           }
//	}
           log.Println("JAPTGET 2 ", err)
          sess.Save(r, w)
            http.Error(w, err.Error(), http.StatusInternalServerError)
	return
 }
// ---------------------------------------------------
// AptGET despliega la pagina del apto
func AptGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := model.Instance(r)
        lisApts, _ := model.Apts()
	v := view.New(r)
	v.Name            = "aparta/apt"
	v.Vars["token"]   = csrfbanana.Token(w, r, sess)
	v.Vars["Title"]   = "Crear Apto"
	v.Vars["Action"]  = "/apto/register"
        v.Vars["LisApts"] =  lisApts
	v.Render(w)
 }
// ---------------------------------------------------
// POST procesa la forma enviada con los datos
func AptPOST(w http.ResponseWriter, r *http.Request) {
        var apt model.Aparta
	sess := model.Instance(r)
        action        := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
	    if validate, missingField := view.Validate(r, []string{"codigo"}); !validate {
                 sess.AddFlash(view.Flash{"Falta Campo: " + missingField, view.FlashError})
                 sess.Save(r, w)
                 AptGET(w, r)
                 return
	      }
              apt.Codigo           = r.FormValue("codigo")
              apt.Descripcion      = r.FormValue("descripcion")
              err := (&apt).AptByCode()

              if err == model.ErrNoResult { // Exito: no hay apartamento creado aun 
                  ex := (&apt).AptCreate()
	          if ex != nil { putErr(ex,sess, "Error Guardando", "AptPOST ") 
	          } else {  // todo bien
                    sess.AddFlash(view.Flash{"Apto. creado: " +apt.Codigo, view.FlashSuccess})
	         }
              }
          }
          sess.Save(r, w)
	  http.Redirect(w, r, "/apto/list", http.StatusFound)
  }

// ---------------------------------------------------
// AptUpGET despliega la pagina del usuario
func AptUpGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
        var apt model.Aparta
	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	id,_ := atoi32(params.ByName("id"))
	SPag   := params.ByName("pg")
        path   :=  fmt.Sprintf("/apto/list/%s", SPag)
        apt.Id = id
        err := (&apt).AptById()
        if err != nil { putErr(err,sess, "No hay aptp", "AptUpGET ") 
            sess.Save(r, w)
            http.Redirect(w, r, path, http.StatusFound)
            return
        }
	v := view.New(r)
	v.Name = "aparta/aptupdate"
	v.Vars["token"]  = csrfbanana.Token(w, r, sess)
	v.Vars["Title"]  = "Actualizar Apto"
	v.Vars["Action"]  = "/apto/update"
        v.Vars["Apto"] = apt
        v.Render(w)
   }
// ---------------------------------------------------
 func   getAptFormUp(a1, a2 model.Aparta, r * http.Request)(stup string){
        var sf string
        var sup []string

        if a1.Descripcion != a2.Descripcion {
	     sf  =  fmt.Sprintf( " descrip = '%s' ", a2.Descripcion )
	     sup = append(sup, sf)
        }
        if a1.Codigo != a2.Codigo {
	     sf  =  fmt.Sprintf( " codigo = '%s' ", a2.Codigo )
	     sup = append(sup, sf)
        }
	lon  := len(sup)
        if  lon > 0 {
            sini        :=  "update apartas set "
	    now         := time.Now()
	    sf           =  fmt.Sprintf( ",  updated_at = '%s' ", now.Format(layout) )
            stup =  strings.Join(sup, ", ")
            sr          :=  fmt.Sprintf(" where apartas.id = %d ", a1.Id)
             stup = sini + stup + sf + sr
        }
        return
  }
// ---------------------------------------------------
// AptUpPOST procesa la forma enviada con los datos
func AptUpPOST(w http.ResponseWriter, r *http.Request) {
        var apt , a2 model.Aparta
	sess := model.Instance(r)
        var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	SId         := params.ByName("id")
        Id,_        := atoi32(SId)
        apt.Id      = Id
	a2.Id       = Id
        path        :=  "/apto/list"
        action        := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
            err                 :=  (&apt).AptById()
            if err != nil { putErr(err,sess, "No hay aptp", "AptUpPOST ")} 
            a2.Codigo           = r.FormValue("codigo")
            a2.Descripcion      = r.FormValue("descripcion")
            st                 :=  getAptFormUp(apt, a2, r)
            if len(st) == 0{
                 sess.AddFlash(view.Flash{"No hay actualizacion solicitada", view.FlashSuccess})
            } else {
             err =  apt.AptUpdate(st)
             if err == nil{
                 sess.AddFlash(view.Flash{"Apto actualizado exitosamente para: " +apt.Codigo, view.FlashSuccess})
             } else    { putErr(err,sess, "Actualizando apto", "AptUpPOST ")     
		sess.Save(r, w)
             }
            }
       }
	http.Redirect(w, r, path, http.StatusFound)
     }
//------------------------------------------------
// AptLisGET displays the aparta page
func AptLisGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
        lisApts, err := model.Apts()
        if err != nil { putErr(err,sess, "Actualizando apto", "AptLisGET ")     
            sess.Save(r, w)
         }
	// Display the view
	v := view.New(r)
	v.Name = "aparta/aptlis"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
        v.Vars["LisApt"]   = lisApts
        v.Vars["Level"]    =  sess.Values["level"]
	v.Render(w)
 }
// ---------------------------------------------------
// AptDeleteGET handles the apto deletion
 func AptDeleteGET(w http.ResponseWriter, r *http.Request) {
        sess := model.Instance(r)
        var apt model.Aparta
        var params httprouter.Params
        params = context.Get(r, "params").(httprouter.Params)
        Id,_ := atoi32(params.ByName("id"))
        apt.Id = Id
        path        :=  "/apto/list"
        err := (&apt).AptById()
        if err != nil { putErr(err,sess, "No tenemos apto", "AptDeleteGET")     
            sess.Save(r, w)
            http.Redirect(w, r, path, http.StatusFound)
            return
        }
	v := view.New(r)
	v.Name = "aparta/aptdelete"
	v.Vars["token"]     = csrfbanana.Token(w, r, sess)
        v.Vars["Title"]     =  "Eliminar Aparta"
        v.Vars["Action"]    =  "/apto/delete"
        v.Vars["Apto"]      =  apt
	v.Render(w)
  }
// ---------------------------------------------------
// AptDeletePOST handles the apto deletion
 func AptDeletePOST(w http.ResponseWriter, r *http.Request) {
        var err error
        sess := model.Instance(r)
        var apt model.Aparta
        var params httprouter.Params
        params = context.Get(r, "params").(httprouter.Params)
        Id,_ := atoi32(params.ByName("id"))
        apt.Id = Id
        path        :=  "/apto/list"
        action        := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
            err  = apt.AptDelete()
            if err != nil { putErr(err,sess, "No posible", "AptDeletePOST")     
            } else {
                sess.AddFlash(view.Flash{"Apto. borrado!", view.FlashSuccess})
            }
            sess.Save(r, w)
        }
	http.Redirect(w, r, path, http.StatusFound)

  }
// ---------------------------------------------------

