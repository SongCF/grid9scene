package model




type Space struct {
	SpaceId    string           // "spaceId123456"
	GridWidth  float32            // 9-grid width
	GridHeight float32            // 9-grid height
	GridM      map[string]*Grid // gridId : Grid
}



func CreateSpace(appId, spaceId string, gridWidth, gridHeight int32) {

}

func DeleteSpace(appId, spaceId string) {

}
