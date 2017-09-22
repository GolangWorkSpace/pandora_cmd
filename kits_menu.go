package main

var _LocalMenu, _Menu *CmdMenu

func SetupMenu() {
	localMenu := &CmdMenu{
		Title: "主菜单",
		SubMenu: []*CmdMenu{
			&CmdMenu{
				Title: "Pod Install加速(本地)",
			},
			&CmdMenu{
				Title: "登录",
				Func:  UserLogin,
			},
			&CmdMenu{
				Title: "注册",
				Func:  UserRegister,
			},
		},
	}
	localMenu.PrepareParentMenu(nil)
	_LocalMenu = localMenu

	menu := &CmdMenu{
		Title: "主菜单",
		SubMenu: []*CmdMenu{
			&CmdMenu{
				Title: "Pod Install加速",
				SubMenu: []*CmdMenu{
					&CmdMenu{
						Title: "使用本地仓库",
					},
					&CmdMenu{
						Title: "使用远端仓库",
					},
				},
			},
			&CmdMenu{
				Title: "团队管理",
				SubMenu: []*CmdMenu{
					&CmdMenu{
						Title: "我的团队",
						Func:  TeamList,
					},
					&CmdMenu{
						Title: "成员管理",
						SubMenu: []*CmdMenu{
							&CmdMenu{
								Title: "成员查询",
								Func:  TeamMembers,
							},
							&CmdMenu{
								Title: "添加成员",
								Func:  TeamMemberAdd,
							},
							&CmdMenu{
								Title: "修改权限",
								Func:  TemMemberModifyRole,
							},
							&CmdMenu{
								Title: "移除成员",
								Func:  TeamMemberRemove,
							},
						},
					},
					&CmdMenu{
						Title: "创建团队",
						Func:  TeamCreate,
					},
				},
			},
			&CmdMenu{
				Title: "项目管理",
				SubMenu: []*CmdMenu{
					&CmdMenu{
						Title: "我的项目",
						Func:  ProjectList,
					},
					&CmdMenu{
						Title: "创建项目",
						Func:  ProjectCreate,
					},
					&CmdMenu{
						Title: "创建项目",
					},
				},
			},
			&CmdMenu{
				Title: "模板管理",
				SubMenu: []*CmdMenu{
					&CmdMenu{
						Title: "模板列表",
						Func:  TemplateList,
					},
					&CmdMenu{
						Title: "查看模板",
						Func:  TemplateShowOne,
					},
					&CmdMenu{
						Title: "添加模板",
					},
					&CmdMenu{
						Title: "跟进版本",
						Func:  TemplateFollow,
					},
					&CmdMenu{
						Title: "修改Pod版本",
					},
					&CmdMenu{
						Title: "添加Pod",
					},
					&CmdMenu{
						Title: "移除Pod",
					},
				},
			},
			&CmdMenu{
				Title: "账户管理",
				SubMenu: []*CmdMenu{
					&CmdMenu{
						Title: "注销",
						Func:  UserLogout,
					},
					&CmdMenu{
						Title: "修改密码",
						Func:  UserChangePassword,
					},
				},
			},
		},
	}
	menu.PrepareParentMenu(nil)
	_Menu = menu
}

func RunParentMenu(aMenu *CmdMenu) {
	if aMenu != nil && aMenu.ParentMenu != nil {
		aMenu.ParentMenu.Run()
	}
}
