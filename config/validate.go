package config

type validate func() error

func (cfg *config) validate() error {
	validations := []validate{
		cfg.db.validate,
	}

	for _, validate := range validations {
		if err := validate(); err != nil {
			return err
		}
	}

	return nil
}
