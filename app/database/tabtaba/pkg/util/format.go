package util

func FormatProjectType(amount float64) (projectType int) {
	if amount <= 3000 {
		projectType = 2 // 3000以下
	} else if amount > 3000 && amount <= 6000 {
		projectType = 3 // 2000-6000
	} else if amount > 6000 && amount <= 12000 {
		projectType = 4 // 5000-12000
	} else if amount > 12000 {
		projectType = 5 // 12000以上
	}
	return projectType
}
