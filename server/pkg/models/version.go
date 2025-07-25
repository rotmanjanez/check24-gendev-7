// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

/*
 * CHECK24 GenDev 7 API
 *
 * API for the 7th CHECK24 GenDev challenge providing product offerings from five different internet providers
 *
 * API version: dev
 */

package models

import (
	"time"
)

// Version - Version information response
type Version struct {

	// The semantic version of the API
	Version string `json:"version,omitempty"`

	// The date and time when the API was built
	BuildDate time.Time `json:"buildDate,omitempty"`

	// The git commit hash of the API build
	CommitHash string `json:"commitHash,omitempty"`
}

// AssertVersionRequired checks if the required fields are not zero-ed
func AssertVersionRequired(obj Version) error {
	return nil
}

// AssertVersionConstraints checks if the values respects the defined constraints
func AssertVersionConstraints(obj Version) error {
	return nil
}
