package rbac

rl_domains = [s_domain |
	some domain_id
	domain := data.domains[domain_id]
	s_domain := {"id": domain_id, "name": domain.name}
]

rl_groups = groups {
	domain_id := input.domain_id	
	all_groups := data.domains[domain_id].groups

	groups := {s_group |
		some group_id
		group := all_groups[group_id]
		s_group := {"id":group_id}
	}
}

rl_subjects = subjects {
	domain_id := input.domain_id	
	all_subjects := data.domains[domain_id].subjects

	subjects := {s_subject |
		some subject_id
		subject := all_subjects[subject_id]
		s_subject := {"id":subject_id, "name": subject.name}		
	}
}

rl_sub_groups = groups {
	domain_id := input.domain_id
	group_pattern := input.group_pattern
	all_groups := data.domains[domain_id].groups

	groups := {group |
		some group
		all_groups[group]
		regex.match(sprintf("^%v$", [group_pattern]), group)
	}
}

rl_domain = result {
	domain_id := input.domain_id
	domain := data.domains[domain_id]
	result := {"id": domain_id, "name": domain.name}
}

rl_permissions_of_subject = output {
	domain_id := input.domain_id
	subject_id := input.subject_id
	subject := data.domains[domain_id].subjects[subject_id]
	domain_groups := data.domains[domain_id].groups
	subject_groups := subject.groups

	domain_permissions := {permission | data.domains[domain_id].permissions[permission]}

	domain_denies := data.domains[domain_id].denies
	subject_groups_denies := {deny_id: deny |
		some group_id, deny_id
		subject_groups[group_id]
		domain_groups[group_id].denies[deny_id]
		deny := domain_denies[deny_id]
	}

	subject_denies := {deny_id: deny |
		some deny_id
		subject.denies[deny_id]
		deny := domain_denies[deny_id]
	}

	denies := object.union(subject_groups_denies, subject_denies)

	denied_permissions := {permission |
		fn_permissions_of_deny(domain_permissions, denies[_])[permission]
	}

	domain_roles := data.domains[domain_id].roles

	subject_groups_roles := {role_id: role |
		some group_id, role_id
		subject_groups[group_id]
		domain_groups[group_id].roles[role_id]
		role := domain_roles[role_id]
	}

	subject_roles := {role_id: role |
		some role_id
		subject.roles[role_id]
		role := domain_roles[role_id]
	}

	roles := object.union(subject_roles, subject_groups_roles)

	permissions_data := [{permission: data} |
		not denied_permissions[permission]
		some role_id
		data := fn_permissions_of_role(domain_permissions, roles[role_id])[permission]
	]

	permissions := {permission |
		some permission
		permissions_data[_][permission]
	}

	output :=  [{"id":permission, "options":data} |
		some permission
		permissions[permission]
		data := [pd |
			pd := permissions_data[_][permission][_]
			pd != {}
		]
	]
}

#     output := {"denies":denies, "denied_permissions":denied_permissions, "roles":roles, "permissions_data":permissions_data, "permissions":permissions, "result": result}

rl_permissions_of_group = output {
	domain_id := input.domain_id
	group_pattern := input.group_pattern
	sub_groups := rl_sub_groups

	domain_groups := data.domains[domain_id].groups

	domain_permissions := {permission | data.domains[domain_id].permissions[permission]}

	domain_denies := data.domains[domain_id].denies
	denies := {deny_id: deny |
		some group_id, deny_id
		sub_groups[group_id]
		domain_groups[group_id].denies[deny_id]
		deny := domain_denies[deny_id]
	}

	denied_permissions := {permission |
		fn_permissions_of_deny(domain_permissions, denies[_])[permission]
	}

	domain_roles := data.domains[domain_id].roles

	roles := {role_id: role |
		some group_id, role_id
		sub_groups[group_id]
		domain_groups[group_id].roles[role_id]
		role := domain_roles[role_id]
	}

	permissions_data := [{permission: data} |
		not denied_permissions[permission]
		some role_id
		data := fn_permissions_of_role(domain_permissions, roles[role_id])[permission]
	]

	permissions := {permission |
		some permission
		permissions_data[_][permission]
	}

	output := [{"id":permission, "options":data} |
		some permission
		permissions[permission]
		data := [pd |
			pd := permissions_data[_][permission][_]
			pd != {}
		]
	]
}

#     output := {"denies":denies, "denied_permissions":denied_permissions, "roles":roles, "permissions_data":permissions_data, "permissions":permissions, "result": result}

rl_permissions_of_role = output {
	domain_id := input.domain_id
	role_id := input.role_id
	role := data.domains[domain_id].roles[role_id]
	domain_permissions := {permission | data.domains[domain_id].permissions[permission]}
	output := fn_permissions_of_role(domain_permissions, role)
}

rl_permissions_of_deny = output {
	domain_id := input.domain_id
	deny_id := input.deny_id
	deny := data.domains[domain_id].denies[deny_id]
	domain_permissions := {permission | data.domains[domain_id].permissions[permission]}
	output := fn_permissions_of_deny(domain_permissions, deny)
}

rl_roles_can_be_accessed_by_subject = output {
	domain_id := input.domain_id
	subject_id := input.subject_id
	all_roles := data.domains[domain_id].roles
	subject_groups := data.domains[domain_id].subjects[subject_id].groups

	output := {{"id":role} |
		some subject_role, subject_group
		subject_groups[subject_group]
		regex.match(sprintf("^%v$", [all_roles[subject_role].scopes.can_be_accessed[_]]), subject_group) == true
		role := subject_role
	}
}

rl_roles_can_be_granted_to_subject = output {
	domain_id := input.domain_id
	subject_id := input.subject_id
	all_roles := data.domains[domain_id].roles
	subject_groups := data.domains[domain_id].subjects[subject_id].groups

	output := {{"id":role} |
		some subject_role, subject_group
		subject_groups[subject_group]
		regex.match(sprintf("^%v$", [all_roles[subject_role].scopes.can_be_granted[_]]), subject_group) == true
		role := subject_role
	}
}

rl_denies_can_be_accessed_by_subject = output {
	domain_id := input.domain_id
	subject_id := input.subject_id
	all_denies := data.domains[domain_id].denies
	subject_groups := data.domains[domain_id].subjects[subject_id].groups

	output := {{"id":deny} |
		some subject_deny, subject_group
		subject_groups[subject_group]
		regex.match(sprintf("^%v$", [all_denies[subject_deny].scopes.can_be_accessed[_]]), subject_group) == true
		deny := subject_deny
	}
}

rl_denies_can_be_granted_to_subject = output {
	domain_id := input.domain_id
	subject_id := input.subject_id
	all_denies := data.domains[domain_id].denies
	subject_groups := data.domains[domain_id].subjects[subject_id].groups

	output := {{"id":deny} |
		some subject_deny, subject_group
		subject_groups[subject_group]
		regex.match(sprintf("^%v$", [all_denies[subject_deny].scopes.can_be_granted[_]]), subject_group) == true
		deny := subject_deny
	}
}

fn_match(pattern, strings) = output {
	contains(pattern, "*.*")
	output := [string |
		string := strings[_]
		regex.match(sprintf("^%v$", [pattern]), string)
	]
}

fn_match(pattern, strings) = output {
	not contains(pattern, "*.*")
	output := [string |
		string := strings[_]
		pattern == string
	]
}

fn_permissions_of_role(domain_permissions, role) = output {
	#获得用户权限模版
	role_included_permission_patterns := {role_included_permission_pattern |
		some role_included_permission_pattern
		role.permissions.included[role_included_permission_pattern]
	}

	#获得用户excluded模版
	role_excluded_permission_patterns := {role_excluded_permission_pattern |
		role_excluded_permission_pattern := role.permissions.excluded[_]
	}

	included_permissions := [{included_permission: pattern} |
		some pattern
		included_permission := fn_match(role_included_permission_patterns[pattern], domain_permissions)[_]
	]

	excluded_permissions := {excluded_permission |
		excluded_permission := fn_match(role_excluded_permission_patterns[_], domain_permissions)[_]
	}

	permissions_data := [{permission: data} |
		some permission
		pattern := included_permissions[_][permission]
		not excluded_permissions[permission]
		data := role.permissions.included[pattern]
	]

	permissions := {permission |
		some permission
		permissions_data[_][permission]
	}

	output := {permission: permission_data |
		some permission
		permissions[permission]
		permission_data := [pd |
			pd := permissions_data[_][permission]
		]
	}
}

fn_permissions_of_deny(domain_permissions, deny) = output {
	#获得用户权限模版
	deny_included_permission_patterns := {deny_included_permission_pattern |
		deny_included_permission_pattern := deny.permissions.included[_]
	}

	#获得用户excluded模版
	deny_excluded_permission_patterns := {deny_excluded_permission_pattern |
		deny_excluded_permission_pattern := deny.permissions.excluded[_]
	}

	included_permissions := {included_permission |
		included_permission := fn_match(deny_included_permission_patterns[_], domain_permissions)[_]
	}

	excluded_permissions := {excluded_permission |
		excluded_permission := fn_match(deny_excluded_permission_patterns[_], domain_permissions)[_]
	}

	output := included_permissions - excluded_permissions
}
