import {JetView, plugins} from "webix-jet";
import {data, serverData} from "models/users";

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
				onselectchange:() => this.aShowForm(),
				"data->onsyncapply":() => this.aShowForm()
			}
		};

		const addButton = { 
			view:"button",
			type:"iconButton",
			label:"Add New User",
			icon:"zmdi zmdi-collection-plus",
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
			{ $subview:"am.empty", name:"content" }
		]};
	}
	urlChange(){
		data.waitData.then(() => {
			const id = this.getParam("userId");
			const list = this.$$("list");
			if (id && list.exists(id)){
				this.$$("list").select(id);
			}
		}); 
	}
	init(){
		this.$$("list").sync(data);
		this.use(plugins.UrlParam, ["userId"]);
	}
	aShowForm(){
		const list = this.$$("list");
		const id = list.getSelectedId();

		if (id === this.aAlreadyShown) return;
		this.aAlreadyShown = id;

		const page = id ? "am.uform" : "am.empty";
		
		this.show({ userId:id },{ silent:true });
		this.show(page, { target:"content" });
	}
}