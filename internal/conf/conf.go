package conf

// при добвлении новой игры CountGames++, в MapToken добавляем данные новой игры
const CountGames = 6

var MapToken = map[int]Token{
	0: {"Riding Extreme 3D", "d28721be-fd2d-4b45-869e-9f253b554e50", "43e35910-c168-4634-ad4f-52fd764a843f"}, //Riding Extreme 3D
	1: {"Chain Cube", "d1690a07-3780-4068-810f-9b5bbf2931b2", "b4170868-cef0-424f-8eb9-be0622e8e8e3"},        //Chain Cube
	2: {"My Clone Army", "74ee0b5b-775e-4bee-974f-63e7f4d5bacb", "fe693b26-b342-4159-8808-15e3ff7f8767"},     //My Clone Army
	3: {"Train Miner", "82647f43-3f87-402d-88dd-09a90025313f", "c4480ac7-e178-4973-8061-9ed5b2e17954"},       //Train Miner
	4: {"Merge Away", "8d1cc2ad-e097-4b86-90ef-7a27e19fb833", "dc128d28-c45b-411c-98ff-ac7726fbaea4"},        //Merge Away
	5: {"Twerk Race", "61308365-9d16-4040-8bb0-2f4a4c69074c", "61308365-9d16-4040-8bb0-2f4a4c69074c"},        //Twerk Race
}

type Token struct {
	GameName string
	AppToken string
	PromoID  string
}
