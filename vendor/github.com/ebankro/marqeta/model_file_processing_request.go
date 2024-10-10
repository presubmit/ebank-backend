/*
 * Marqeta Core API
 *
 * Simplified management of your payment programs
 *
 * API version: 3.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package marqeta

type FileProcessingRequest struct {
	Token           string `json:"token,omitempty"`
	FileProcessType string `json:"file_process_type"`
	SourceFile      string `json:"source_file"`
	ArchiveFile     string `json:"archive_file,omitempty"`
}
