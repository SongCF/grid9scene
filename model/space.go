package model




type Space struct {
	SpaceId string		     // "spaceId123456"
	GridWidth int32		     // 9-grid width
	GridHeight int32	     // 9-grid height
	GridList map[string]*Grid    // gridId : Grid
}



func CreateSpace(appId, spaceId string, gridWidth, gridHeight int32) {

}

func DeleteSpace(appId, spaceId string) {

}
