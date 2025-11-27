package model

type NSPDResp struct {
	Data struct {
		Type     string `json:"type"`
		Features []struct {
			Id       int `json:"id"`
			Geometry struct {
				Type        string        `json:"type"`
				Coordinates [][][]float64 `json:"coordinates"`
				Crs         struct {
					Type       string `json:"type"`
					Properties struct {
						Name string `json:"name"`
					} `json:"properties"`
				} `json:"crs"`
			} `json:"geometry"`
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

type TerrZone struct {
	Type     string `json:"type"`
	Features []struct {
		Id       int    `json:"id"`
		Type     string `json:"type"`
		Geometry struct {
			Type        string        `json:"type"`
			Coordinates [][][]float64 `json:"coordinates"`
			Crs         struct {
				Type       string `json:"type"`
				Properties struct {
					Name string `json:"name"`
				} `json:"properties"`
			} `json:"crs"`
		} `json:"geometry"`
		Properties struct {
			CadastralDistrictsCode int    `json:"cadastralDistrictsCode"`
			Category               int    `json:"category"`
			CategoryName           string `json:"categoryName"`
			Descr                  string `json:"descr"`
			ExternalKey            string `json:"externalKey"`
			InteractionId          int    `json:"interactionId"`
			Label                  string `json:"label"`
			Options                struct {
				CadastralDistrict           string `json:"cadastral_district"`
				ContentRestrictEncumbrances string `json:"content_restrict_encumbrances"`
				LegalActDocumentDate        string `json:"legal_act_document_date"`
				LegalActDocumentIssuer      string `json:"legal_act_document_issuer"`
				LegalActDocumentName        string `json:"legal_act_document_name"`
				LegalActDocumentNumber      string `json:"legal_act_document_number"`
				NameByDoc                   string `json:"name_by_doc"` //TODO: use this
				OldAccountNumber            string `json:"old_account_number"`
				PermittedUsesName           string `json:"permitted_uses_name"`
				RegNumbBorder               string `json:"reg_numb_border"`
				RegistrationDate            string `json:"registration_date"`
				TypeBoundaryValue           string `json:"type_boundary_value"`
				TypeZone                    string `json:"type_zone"`
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
}
