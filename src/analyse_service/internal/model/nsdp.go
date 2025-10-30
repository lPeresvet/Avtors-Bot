package model

type NSPDResp struct {
	Data struct {
		Type     string `json:"type"`
		Features []struct {
			Id         int `json:"id"`
			Properties struct {
				CadastralDistrictsCode int    `json:"cadastralDistrictsCode"`
				Category               int    `json:"category"`
				CategoryName           string `json:"categoryName"`
				Descr                  string `json:"descr"`
				ExternalKey            string `json:"externalKey"`
				InteractionId          int    `json:"interactionId"`
				Label                  string `json:"label"`
				Options                struct {
					Area                              interface{} `json:"area"`
					CadNum                            string      `json:"cad_num"`
					CommonDataStatus                  string      `json:"common_data_status"`
					CostApplicationDate               string      `json:"cost_application_date"`
					CostApprovementDate               string      `json:"cost_approvement_date"`
					CostDeterminationDate             string      `json:"cost_determination_date"`
					CostIndex                         float64     `json:"cost_index"`
					CostRegistrationDate              string      `json:"cost_registration_date"`
					CostValue                         float64     `json:"cost_value"`
					DeclaredArea                      interface{} `json:"declared_area"`
					DeterminationCouse                string      `json:"determination_couse"`
					LandRecordArea                    interface{} `json:"land_record_area"`
					LandRecordAreaDeclaration         interface{} `json:"land_record_area_declaration"`
					LandRecordAreaVerified            int         `json:"land_record_area_verified"`
					LandRecordCategoryType            string      `json:"land_record_category_type"`
					LandRecordRegDate                 string      `json:"land_record_reg_date"`
					LandRecordSubtype                 string      `json:"land_record_subtype"`
					LandRecordType                    string      `json:"land_record_type"`
					OwnershipType                     string      `json:"ownership_type"`
					PermittedUseEstablishedByDocument string      `json:"permitted_use_established_by_document"`
					PreviouslyPosted                  string      `json:"previously_posted"`
					QuarterCadNumber                  string      `json:"quarter_cad_number"`
					ReadableAddress                   string      `json:"readable_address"`
					RegistrationDate                  string      `json:"registration_date"`
					RightType                         string      `json:"right_type"`
					SpecifiedArea                     int         `json:"specified_area"`
					Status                            string      `json:"status"`
					Subtype                           string      `json:"subtype"`
				} `json:"options"`
				Subcategory int `json:"subcategory"`
				SystemInfo  struct {
					Inserted   string `json:"inserted"`
					InsertedBy string `json:"insertedBy"`
					Updated    string `json:"updated"`
					UpdatedBy  string `json:"updatedBy"`
				} `json:"systemInfo"`
			} `json:"properties"`
		} `json:"features"`
	} `json:"data"`
	//Meta []struct {
	//	TotalCount int `json:"totalCount"`
	//	CategoryId int `json:"categoryId"`
	//} `json:"meta"`
}
