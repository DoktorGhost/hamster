package conf

// при добвлении новой игры CountGames++, в MapToken добавляем данные новой игры
const CountGames = 17

var MapToken = map[int]Token{
	0:  {"Zoopolis", "b2436c89-e0aa-4aed-8046-9b0515e1c46b", "b2436c89-e0aa-4aed-8046-9b0515e1c46b"},
	1:  {"Mow and Trim", "ef319a80-949a-492e-8ee0-424fb5fc20a6", "ef319a80-949a-492e-8ee0-424fb5fc20a6"},
	2:  {"Chain Cube", "d1690a07-3780-4068-810f-9b5bbf2931b2", "b4170868-cef0-424f-8eb9-be0622e8e8e3"},
	3:  {"Train Miner", "82647f43-3f87-402d-88dd-09a90025313f", "c4480ac7-e178-4973-8061-9ed5b2e17954"},
	4:  {"Merge Away", "8d1cc2ad-e097-4b86-90ef-7a27e19fb833", "dc128d28-c45b-411c-98ff-ac7726fbaea4"},
	5:  {"Twerk Race", "61308365-9d16-4040-8bb0-2f4a4c69074c", "61308365-9d16-4040-8bb0-2f4a4c69074c"},
	6:  {"Polysphere", "2aaf5aee-2cbc-47ec-8a3f-0962cc14bc71", "2aaf5aee-2cbc-47ec-8a3f-0962cc14bc71"},
	7:  {"Tile Trio", "e68b39d2-4880-4a31-b3aa-0393e7df10c7", "e68b39d2-4880-4a31-b3aa-0393e7df10c7"},
	8:  {"Stone Age", "04ebd6de-69b7-43d1-9c4b-04a6ca3305af", "04ebd6de-69b7-43d1-9c4b-04a6ca3305af"},
	9:  {"Fluff Crusade", "112887b0-a8af-4eb2-ac63-d82df78283d9", "112887b0-a8af-4eb2-ac63-d82df78283d9"},
	10: {"Bouncemasters", "bc72d3b9-8e91-4884-9c33-f72482f0db37", "bc72d3b9-8e91-4884-9c33-f72482f0db37"},
	11: {"Hide Ball", "4bf4966c-4d22-439b-8ff2-dc5ebca1a600", "4bf4966c-4d22-439b-8ff2-dc5ebca1a600"},
	12: {"Pin Out Master", "d2378baf-d617-417a-9d99-d685824335f0", "d2378baf-d617-417a-9d99-d685824335f0"},
	13: {"Count Masters", "4bdc17da-2601-449b-948e-f8c7bd376553", "4bdc17da-2601-449b-948e-f8c7bd376553"},
	14: {"Infected Frontier", "eb518c4b-e448-4065-9d33-06f3039f0fcb", "eb518c4b-e448-4065-9d33-06f3039f0fcb"},
	15: {"Among Waterr", "daab8f83-8ea2-4ad0-8dd5-d33363129640", "daab8f83-8ea2-4ad0-8dd5-d33363129640"},
	16: {"Factory World", "d02fc404-8985-4305-87d8-32bd4e66bb16", "d02fc404-8985-4305-87d8-32bd4e66bb16"},
}

type Token struct {
	GameName string
	AppToken string
	PromoID  string
}
