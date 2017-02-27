package routing

import "github.com/hoopra/GoAuthServer/controllers"

// Port on which this router should listen
var Port = ":8000"

// Routes available on this router
var Routes = []route{
	route{"Register", "POST", "/register", controllers.Register, false},
	route{"Login", "POST", "/login", controllers.Login, false},
	route{"Refresh", "POST", "/refresh", controllers.RefreshToken, true},
	route{"Logout", "GET", "/logout", controllers.Logout, true},

	route{"Update Name", "POST", "/update/name", controllers.UpdateName, true},
	route{"Update Password", "POST", "/update/password", controllers.UpdatePassword, true},
}
