import {JetView} from "webix-jet";
import {data, serverData} from "models/users";

import form from "./uform";

export default class TopView extends JetView{
	config(){
		const user = this.app.getService("user");
		user.guardRight(user.rights.AdminUser);

		const list = {
			view:"list", id:"list",
			scroll:"auto", width:200, select:true,
			type:{
				height:70, template:"#name#<br>#email#"
			},
			on:{
				onAfterSelect:id => this.show({ id })
			}
		};

		const addButton = { 
			view:"button",
			type:"iconButton",
			label:"Add New User",
			icon:"zmdi zmdi-plus",
			width:200, click: () => {
				serverData.save(
					0,
					"insert",
					{ name:"New User", email:"", rights:"" }
				).then(obj => {
					data.add(obj);
					this.$$("list").select(obj.id);
				});
			}
		};

		return { type:"space", cols:[ 
			{ type:"clean", rows:[ list, addButton ] },
			form
		]};
	}
	urlChange(){
		data.waitData.then(() => {
			const id = this.getParam("id");
			if (id){
				this.$$("list").select(id);
			} 
		}); 
	}
	init(){
		this.$$("list").sync(data);
	}
}