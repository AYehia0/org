package controllers

// here we put all the request and response types for the organization controller

// Creating a new organization request
type CreateOrganizationRequest struct {
	Name string `json:"name" binding:"required"`
	Desc string `json:"description" binding:"required"`
}

// Updating an organization request
type UpdateOrganizationRequest struct {
	Name string `json:"name"`
	Desc string `json:"description"`
}

// Invite a user to an organization request
type InviteUserToOrganizationRequest struct {
	Email string `json:"user_email" binding:"required,email"`
}
