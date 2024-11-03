package dto

func MapMany[TSource any, TDest any](source []*TSource, mapper func(*TSource) *TDest) []*TDest {
	result := make([]*TDest, len(source))
	for i, tSource := range source {
		result[i] = mapper(tSource)
	}
	return result
}
