export const data = new webix.DataCollection({ 
	url: () => remote.api.admin.GetUsers()
});

export const serverData = {
	save: remote.api.admin.SaveUser
};
