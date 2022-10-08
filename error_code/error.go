package error_code

type ErrorCode struct {
	Err         error
	Status      int
	Description string
}

func (c ErrorCode) CheckValid() bool {
	if c.Err != nil {
		return false
	} else {
		return true
	}
}
