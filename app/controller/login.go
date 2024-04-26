package controller

import (
	"fmt"
	"log"
	"net/http"
        "encoding/json"
        "encoding/base64"

	"github.com/gatopardo/hcondo/app/model"
	"github.com/gatopardo/hcondo/app/shared/passhash"
	"github.com/gatopardo/hcondo/app/shared/view"

	"github.com/gorilla/sessions"
	"github.com/josephspurrier/csrfbanana"

        "github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"

  )

const (
	// Name of the session variable that tracks login attempts
	sessLoginAttempt = "login_attempt"
  )

// loginAttempt increments the number of login attempts in sessions variable
func loginAttempt(sess *sessions.Session) {
	// Log the attempt
	if sess.Values[sessLoginAttempt] == nil {
		sess.Values[sessLoginAttempt] = 1
	} else {
		sess.Values[sessLoginAttempt] = sess.Values[sessLoginAttempt].(int) + 1
	}
  }

 func PadRight(str, pad string, lenght int) string {
	for {
		str += pad
		if len(str) > lenght {
			return str[0:lenght]
		}
	}
  }

// LoginGET displays the login page
func LoginGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
	v := view.New(r)
	v.Name = "login/login"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	// Refill any form fields
	view.Repopulate([]string{"cuenta","password"}, r.Form, v.Vars)
	v.Render(w)
  }

// LoginPOST handles the login form submission
 func LoginPOST(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
	// Validate with required fields
	if validate, missingField := view.Validate(r, []string{"cuenta", "pass"}); !validate {
		sess.AddFlash(view.Flash{"Falta campo: " + missingField, view.FlashError})
		sess.Save(r, w)
		LoginGET(w, r)
		return
	}
       // Prevenir intentos login de fuerza bruta pretendiendo entrada invalida :-)
        if sess.Values[sessLoginAttempt] != nil && sess.Values[sessLoginAttempt].(int) >= 3 {
		log.Println("Intentos de Entrada Repetidos en Exceso")
		sess.AddFlash(view.Flash{"No mas intentos :-)", view.FlashNotice})
		sess.Save(r, w)
		LoginGET(w, r)
		return
	}
	// Form values
	cuenta := r.FormValue("cuenta")
	password := r.FormValue("pass")
	// Get database user
         var user  model.User
         user.Cuenta  = cuenta
	 err := (&user).UserByCuenta()
	if err == model.ErrNoResult {
		loginAttempt(sess)
		sess.AddFlash(view.Flash{"Cuenta incorrecta - Intento: " + fmt.Sprintf("%v", sess.Values[sessLoginAttempt]), view.FlashWarning})
                log.Println("Cuenta Incorreta ", err)
//                fmt.Println("Cuenta Incorreta ", err)
	} else if err != nil {
		log.Println(" Error busqueda ",err)
		sess.AddFlash(view.Flash{"Un error. Favor probar mas tarde.", view.FlashError})
		sess.Save(r, w)
	} else if passhash.MatchString(user.Password, password) {
             if user.Nivel == 0 {
                 sess.AddFlash(view.Flash{"Cuenta inactiva entrada prohibida.", view.FlashNotice})
//			sess.Save(r, w)
		} else { // Login successfully
			model.Empty(sess)
			sess.AddFlash(view.Flash{"Entrada exitosa!", view.FlashSuccess})
			sess.Values["id"] = user.Id
			sess.Values["cuenta"] = cuenta
			sess.Values["level"] = user.Nivel
			sess.Save(r, w)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	} else {
		loginAttempt(sess)
		sess.AddFlash(view.Flash{"Clave incorrecta - Intento: " + fmt.Sprintf("%v", sess.Values[sessLoginAttempt]), view.FlashWarning})
		sess.Save(r, w)
	}
	// Show the login page again
	LoginGET(w, r)
  }
  
// JLoginGET service to return persons data
 func JLoginGET(w http.ResponseWriter, r *http.Request) {
// fmt.Printf("JLoginGET 0 \n")
        var params httprouter.Params
	sess           := model.Instance(r)

	v := view.New(r)
	v.Vars["token"] = csrfbanana.Token(w, r, sess)

        params          = context.Get(r, "params").(httprouter.Params)
        cuenta         := params.ByName("cuenta")
        password       := params.ByName("password")
// fmt.Printf("JLoginGET 1 %s %s\n", cuenta, password)
        stEnc, _ := base64.StdEncoding.DecodeString(password)
	   password = string(stEnc)
// fmt.Printf("JLoginGET 2 %s %s\n", cuenta, password)
        var jpers  model.Jperson
        jpers.Cuenta  = cuenta
	pass, err    := (&jpers).JPersByCuenta()
	if err == model.ErrNoResult {
             loginAttempt(sess)
	} else {
// fmt.Printf("JLoginGET 3 %s %s\n", pass, password)
		b:= passhash.MatchString(pass, password)
// fmt.Printf("JLoginGET 4 %t\n", b)
                if b && jpers.Nivel > 0{
		   var js []byte
		   js, err =  json.Marshal(jpers)
// fmt.Printf("JLoginGET 5 %s\n", string(js))
                   if err == nil{
			model.Empty(sess)
			sess.Values["id"] = jpers.Id
                        sess.Save(r, w)
                        w.Header().Set("Content-Type", "application/json")
                        w.Write(js)
			return
                    }
	        }
            }
               log.Println(err)
//	       http.Error(w, err.Error(), http.StatusBadRequest)
               http.Error(w, err.Error(), http.StatusInternalServerError)
		return
      }



// LogoutGET clears the session and logs the user out
func LogoutGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := model.Instance(r)

	// If user is authenticated
	if sess.Values["id"] != nil {
	        model.Empty(sess)
		sess.AddFlash(view.Flash{"Goodbye!", view.FlashNotice})
		sess.Save(r, w)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
