import {JetView} from "webix-jet";
import {data, serverData} from "models/users";

export default class TopView extends JetView{
	config(){
		const user = this.app.getService("user");
		user.guardRight(user.rights.AdminUser);

		const buttons = {
			width: 120, rows:[
				{ view:"button", value:"Save", align:"right", type:"form", click:() => {
					const form = this.$$("form");
					const id = form.getValues().id;

					data.updateItem(id, form.getValues());
					serverData.save(id, "update", form.getValues());
				}},
				{ view:"button", value:"Delete", click:() => {
					const form = this.$$("form");
					const id = form.getValues().id;
					webix.confirm(`This will disable access for ${data.getItem(id).name}`).then(() => {
						serverData.save(id, "delete", {});
						data.remove(id);
					});
				}}
			]
		};

		const rights = [];
		for (var key in user.rights){
			rights.push({id:""+user.rights[key], value:key});
		}
		const rightsSelector = {
			view:"dbllist", maxWidth:600,
			list:{ yCount:4 },
			name:"rights",
			labelLeft:"Denied to",
			labelRight:"Allowed to",
			data:rights
		};

		const editors = {
			gravity: 9999, maxWidth: 400,
			rows:[
				{ view:"text", name:"name", label: "Name" },
				{ view:"text", name:"email", label: "Email" },
			]
		};


		return {
			view:"form", id:"form", maxWidth:600, rows:[
				{ cols:[editors, {}, buttons ] },
				rightsSelector,
				{}
			]
		};
	}
	urlChange(){
		data.waitData.then(() => {
			const id = this.getParam("userId", true);
			if (id){
				this.$$("form").setValues( data.getItem(id) );
				this.$$("form").focus();
			}
		}); 
	}
}