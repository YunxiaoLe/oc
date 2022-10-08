package logit

// LogWithInfo 用于作为log文件通知
func LogWithInfo(title string, infos map[string]interface{}) {
	log.WithFields(infos).Info(title)
}
func LogWithError(title string, infos map[string]interface{}) {
	log.WithFields(infos).Error(title)
}
func LogWithWarn(title string, infos map[string]interface{}) {
	log.WithFields(infos).Warn(title)
}
func LogWithFatal(title string, infos map[string]interface{}) {
	log.WithFields(infos).Fatal(title)
}
func LogWithPanic(title string, infos map[string]interface{}) {
	log.WithFields(infos).Panic(title)
}
