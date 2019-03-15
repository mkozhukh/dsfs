import "./styles/app.css";
import {JetApp} from "webix-jet";

webix.ready(() => {
	var app = new JetApp({
		id:			APPNAME,
		version:	VERSION,
		start:		"/am.top/am.users",
		debug:		!PRODUCTION
	});
	app.render();

	app.attachEvent("app:error:resolve", function(name, error){
		window.console.error(error);
	});

	const access = {
		rights: remote.data.rights,
		getUser: () => remote.data.user,
		getRights: () => access.getUser().rights,
		hasRight: right => access.getRights().indexOf(right) !== -1,
		guardRight: right => {
			if (!access.hasRight(right)) throw "Access denied";
		}
	};

	app.setService("user", access);
});