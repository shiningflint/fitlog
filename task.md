golang net/http

GET	/admin/users		list all users
GET	/admin/users/new	show new user form
POST	/admin/users		create new user, redirect to edit page
GET	/admin/users/:id/edit	show user edit form
PUT	/admin/users/:id	update user information
PATCH	/admin/users/:id	update user information
DELETE	/admin/users/:id	destroy user information
