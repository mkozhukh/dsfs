import {JetView} from "webix-jet";

export default class TopView extends JetView{
	config(){
		const user = this.app.getService("user");
		user.guardRight(user.rights.AdminUser);

		const toolbar = {
			view:"toolbar", css:"webix_dark", padding:{ left: 10, right:10 }, elements:[
				{ view:"label", label:"Access management" }
				//{ view:"segmented", options:["Users", "Roles"], width: 320 }
			]
		};

		return {
			rows:[
				toolbar,
				{ $subview:true }
			]
		};
	}
	init(){
		const root = this.getRoot();
		webix.extend(root, webix.ProgressBar);

		this.on(webix, "onRemoteError", err => {
			window.console.error(err);
			webix.message({ type:"error", text:"Server side error<br>"+err });
		});

		this.on(webix, "onRemoteCall", (res) => {
			root.showProgress({ type:"top", delay: 2000 });
			res.finally(() => {
				root.hideProgress();
			});
		});
	}
}