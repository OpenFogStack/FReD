package app

// HandleCreateKeygroup creates a new keygroup with the specified name in Storage.
func (a *App) HandleCreateKeygroup(kgname string) error {
	if err := checkParameters(kgname); err != nil {
		return err
	}

	return a.createKeygroup(kgname)
}

// HandleDeleteKeygroup removes the keygroup with the specified name from Storage.
func (a *App) HandleDeleteKeygroup(kgname string) error {
	if err := checkParameters(kgname); err != nil {
		return err
	}

	return a.deleteKeygroup(kgname)

	return nil
}