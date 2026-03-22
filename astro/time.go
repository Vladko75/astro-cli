package main

import "time"

// Location holds coordinates, description, and timezone for a place
type Location struct {
	Lat         float64
	Lon         float64
	Description string
	Timezone    string // IANA timezone string (e.g., "America/New_York")
}

// cities maps city names to their coordinates, descriptions, and timezones
var cities = map[string]Location{
	"jakarta":       {-6.2088, 106.8456, "Capital of Indonesia, largest city in Southeast Asia", "Asia/Jakarta"},
	"dhaka":         {23.8103, 90.4125, "Capital of Bangladesh, densely populated megacity", "Asia/Dhaka"},
	"tokyo":         {35.6762, 139.6503, "Capital of Japan, largest metropolitan area in the world", "Asia/Tokyo"},
	"delhi":         {28.7041, 77.1025, "Capital of India, center of Indian culture and history", "Asia/Kolkata"},
	"shanghai":      {31.2304, 121.4737, "Largest city in China, major financial hub", "Asia/Shanghai"},
	"guangzhou":     {23.1291, 113.2644, "Major city in southern China, important trade center", "Asia/Shanghai"},
	"cairo":         {30.0444, 31.2357, "Capital of Egypt, largest city in Africa", "Africa/Cairo"},
	"manila":        {14.5995, 120.9842, "Capital of the Philippines, major port city", "Asia/Manila"},
	"kolkata":       {22.5726, 88.3639, "Major city in India, cultural and intellectual center", "Asia/Kolkata"},
	"seoul":         {37.5665, 126.9780, "Capital of South Korea, vibrant modern metropolis", "Asia/Seoul"},
	"karachi":       {24.8607, 67.0011, "Largest city in Pakistan, major seaport", "Asia/Karachi"},
	"mumbai":        {19.0760, 72.8777, "Financial capital of India, entertainment hub", "Asia/Kolkata"},
	"sao paulo":     {-23.5505, -46.6333, "Largest city in Brazil, major financial center", "America/Sao_Paulo"},
	"bangkok":       {13.7563, 100.5018, "Capital of Thailand, bustling cultural city", "Asia/Bangkok"},
	"mexico city":   {19.4326, -99.1332, "Capital of Mexico, built on the site of Aztec city", "America/Mexico_City"},
	"beijing":       {39.9042, 116.4074, "Capital of China, political and cultural center", "Asia/Shanghai"},
	"lahore":        {31.5497, 74.3436, "Second largest city in Pakistan, cultural heart", "Asia/Karachi"},
	"istanbul":      {41.0082, 28.9784, "Bridge between Europe and Asia, historic city", "Europe/Istanbul"},
	"moscow":        {55.7558, 37.6173, "Capital of Russia, largest city in Europe", "Europe/Moscow"},
	"ho chi minh city": {10.7769, 106.7009, "Largest city in Vietnam, economic hub", "Asia/Ho_Chi_Minh"},
	"buenos aires":  {-34.6037, -58.3816, "Capital of Argentina, cultural metropolis", "America/Argentina/Buenos_Aires"},
	"new york city": {40.7128, -74.0060, "Largest city in USA, global financial center", "America/New_York"},
	"shenzhen":      {22.5431, 114.0579, "Modern city in China, special economic zone", "Asia/Shanghai"},
	"bengaluru":     {12.9716, 77.5946, "IT capital of India, tech hub", "Asia/Kolkata"},
	"osaka":         {34.6937, 135.5023, "Major city in Japan, center of commerce", "Asia/Tokyo"},
	"lagos":         {6.5244, 3.3792, "Largest city in Africa by population", "Africa/Lagos"},
	"los angeles":   {34.0522, -118.2437, "Major city in USA, entertainment center", "America/Los_Angeles"},
	"chennai":       {13.0827, 80.2707, "Major city in India, cultural capital of South", "Asia/Kolkata"},
	"kinshasa":      {-4.2613, 15.3136, "Capital of Democratic Republic of Congo", "Africa/Kinshasa"},
	"bogota":        {4.7110, -74.0721, "Capital of Colombia, high-altitude metropolis", "America/Bogota"},
	"lima":          {-12.0464, -77.0428, "Capital of Peru, gateway to South America", "America/Lima"},
	"london":        {51.5074, -0.1278, "Capital of United Kingdom, global cultural center", "Europe/London"},
	"rio de janeiro": {-22.9068, -43.1729, "Major city in Brazil, famous for beaches and Christ statue", "America/Sao_Paulo"},
	"paris":         {48.8566, 2.3522, "Capital of France, city of light and romance", "Europe/Paris"},
	"hyderabad":     {17.3850, 78.4867, "Major city in India, IT and pearl capital", "Asia/Kolkata"},
	"tehran":        {35.6892, 51.3890, "Capital of Iran, largest city in Middle East", "Asia/Tehran"},
	"luanda":        {-8.8383, 13.2344, "Capital of Angola, major port city", "Africa/Luanda"},
	"bandung":       {-6.9175, 107.6191, "City in Indonesia, center for textiles", "Asia/Jakarta"},
	"kuala lumpur":  {3.1390, 101.6869, "Capital of Malaysia, modern metropolis", "Asia/Kuala_Lumpur"},
	"dar es salaam": {-6.8000, 39.2833, "Largest city in Tanzania, major port", "Africa/Dar_es_Salaam"},
	"suzhou":        {31.2989, 120.5954, "Ancient city in China, famous for gardens", "Asia/Shanghai"},
	"ahmedabad":     {23.0225, 72.5714, "Major city in India, textile center", "Asia/Kolkata"},
	"hangzhou":      {30.2875, 120.1551, "City in China, historic silk trade hub", "Asia/Shanghai"},
	"wuhan":         {30.5928, 114.3055, "Major city in China, transport hub", "Asia/Shanghai"},
	"tianjin":       {39.0842, 117.2010, "Major seaport in China", "Asia/Shanghai"},
	"alexandria":    {31.2001, 29.9187, "Historic city in Egypt, founded by Alexander the Great", "Africa/Cairo"},
	"nagoya":        {35.1815, 136.9066, "Major city in Japan, industrial center", "Asia/Tokyo"},
	"johannesburg":  {-26.2023, 28.0436, "Largest city in South Africa, economic hub", "Africa/Johannesburg"},
	"chongqing":     {29.5630, 106.5516, "Major city in China, mountain metropolis", "Asia/Shanghai"},
	"riyadh":        {24.7136, 46.6753, "Capital of Saudi Arabia, modern desert city", "Asia/Riyadh"},
	"surat":         {21.1702, 72.8311, "Major city in India, diamond and textile center", "Asia/Kolkata"},
	"surabaya":      {-7.2506, 112.7508, "Major city in Indonesia, maritime hub", "Asia/Jakarta"},
	"pune":          {18.5204, 73.8567, "Major city in India, educational center", "Asia/Kolkata"},
	"khartoum":      {15.5007, 32.5599, "Capital of Sudan, meeting point of Nile tributaries", "Africa/Khartoum"},
	"nanjing":       {32.0603, 118.7969, "Historic city in China, former capital", "Asia/Shanghai"},
	"santiago":      {-33.8688, -70.2093, "Capital of Chile, surrounded by Andes", "America/Santiago"},
	"chicago":       {41.8781, -87.6298, "Major city in USA, architectural hub", "America/Chicago"},
	"chengdu":       {30.5728, 104.0668, "Major city in China, panda conservation center", "Asia/Shanghai"},
	"xi'an":         {34.3416, 108.9398, "Ancient capital of China, Silk Road gateway", "Asia/Shanghai"},
	"hong kong":     {22.2793, 114.1745, "Global financial center, vibrant harbor city", "Asia/Hong_Kong"},
	"dongguan":      {23.0489, 113.7381, "Major manufacturing city in China", "Asia/Shanghai"},
	"foshan":        {23.0298, 113.1210, "Major city in China, ceramic production hub", "Asia/Shanghai"},
	"shenyang":      {41.8045, 123.4328, "Major city in northeastern China", "Asia/Shanghai"},
	"baghdad":       {33.3128, 44.3615, "Capital of Iraq, historic Mesopotamian city", "Asia/Baghdad"},
	"madrid":        {40.4168, -3.7038, "Capital of Spain, cultural capital", "Europe/Madrid"},
	"harbin":        {45.8038, 126.5340, "Northern city in China, ice sculpture hub", "Asia/Shanghai"},
	"houston":       {29.7604, -95.3698, "Major city in USA, energy industry center", "America/Chicago"},
	"dallas":        {32.7767, -96.7970, "Major city in USA, business and arts center", "America/Chicago"},
	"toronto":       {43.6532, -79.3832, "Largest city in Canada, multicultural metropolis", "America/Toronto"},
	"miami":         {25.7617, -80.1918, "Major city in USA, tropical gateway", "America/New_York"},
	"belo horizonte": {-19.9191, -43.9386, "Major city in Brazil, modernist urban center", "America/Sao_Paulo"},
	"singapore":     {1.3521, 103.8198, "City-state, major global financial hub", "Asia/Singapore"},
	"philadelphia":  {39.9526, -75.1652, "Major city in USA, historic Liberty city", "America/New_York"},
	"atlanta":       {33.7490, -84.3880, "Major city in USA, hub of the New South", "America/New_York"},
	"fukuoka":       {33.5904, 130.4017, "Major city in Japan, gateway to Asia", "Asia/Tokyo"},
	"barcelona":     {41.3851, 2.1734, "Major city in Spain, architectural marvel", "Europe/Madrid"},
	"saint petersburg": {59.9311, 30.3609, "Historic city in Russia, Venice of the North", "Europe/Moscow"},
	"qingdao":       {36.0671, 120.3826, "Coastal city in China, beer capital", "Asia/Shanghai"},
	"dalian":        {38.9140, 121.6147, "Coastal city in China, port metropolis", "Asia/Shanghai"},
	"washington d.c.": {38.9072, -77.0369, "Capital of USA, political and cultural center", "America/New_York"},
	"yangon":        {16.8661, 96.1951, "Largest city in Myanmar, golden pagoda hub", "Asia/Yangon"},
	"jinan":         {36.6519, 117.1205, "Major city in China, spring city", "Asia/Shanghai"},
	"guadalajara":   {20.6595, -103.2494, "Major city in Mexico, birthplace of mariachi", "America/Mexico_City"},
	"izmir":         {38.4237, 27.1428, "Coastal city in Turkey, Aegean gateway", "Europe/Istanbul"},
	"vancouver":     {49.2827, -123.1207, "Major city in Canada, Pacific urban beauty", "America/Vancouver"},
	"incheon":       {37.2757, 126.7335, "Major city in South Korea, gateway city", "Asia/Seoul"},
	"almaty":        {43.2389, 76.9453, "Largest city in Kazakhstan, mountain surrounded", "Asia/Almaty"},
	"muscat":        {23.6100, 58.5400, "Capital of Oman, port city on Arabian Sea", "Asia/Muscat"},
	"phoenix":       {33.4484, -112.0742, "Major city in USA, desert oasis", "America/Phoenix"},
	"antananarivo":  {-18.8792, 47.5079, "Capital of Madagascar, highland city", "Africa/Nairobi"},
	"casablanca":    {33.5731, -7.5898, "Major city in Morocco, economic hub", "Africa/Casablanca"},
	"jeddah":        {21.5433, 39.1727, "Major city in Saudi Arabia, port on Red Sea", "Asia/Riyadh"},
	"semarang":      {-6.9667, 110.4167, "Major city in Indonesia, central Java hub", "Asia/Jakarta"},
	"medan":         {2.1956, 98.6722, "Major city in Indonesia, gateway to Sumatra", "Asia/Jakarta"},
	"palembang":     {-2.8949, 104.7533, "Major city in Indonesia, port metropolis", "Asia/Jakarta"},
	"quezon city":   {14.6349, 121.0388, "Major city in Philippines, administrative center", "Asia/Manila"},
	"cebu":          {10.3157, 123.8854, "Major city in Philippines, Visayan hub", "Asia/Manila"},
	"davao":         {7.1907, 125.4553, "Major city in Philippines, Mindanao metropolis", "Asia/Manila"},
	"ayutthaya":     {14.3569, 100.5846, "Historic city in Thailand, ancient capital", "Asia/Bangkok"},
	"hanoi":         {21.0285, 105.8542, "Capital of Vietnam, millenary capital", "Asia/Ho_Chi_Minh"},
	"haiphong":      {20.8554, 106.6843, "Major city in Vietnam, northern port", "Asia/Ho_Chi_Minh"},
	"can tho":       {10.0379, 105.7869, "Major city in Vietnam, Mekong delta hub", "Asia/Ho_Chi_Minh"},
	"phnom penh":    {11.5564, 104.9282, "Capital of Cambodia, riverside metropolis", "Asia/Bangkok"},
	"siem reap":     {13.3670, 103.8452, "City in Cambodia, gateway to Angkor temples", "Asia/Bangkok"},
	"mandalay":      {21.9588, 96.0853, "Major city in Myanmar, cultural center", "Asia/Yangon"},
	"naypyidaw":     {19.8183, 96.1921, "Capital of Myanmar, planned city", "Asia/Yangon"},
	// Polar regions
	"north pole":    {90.0, 0.0, "Geographic North Pole, northernmost point on Earth", "UTC"},
	"south pole":    {-90.0, 0.0, "Geographic South Pole, southernmost point on Earth", "UTC"},
	"mcmurdo station": {-77.85, 166.67, "Antarctic research station, US base", "UTC"},
	"alert":         {82.5, -62.3, "Northernmost settlement in Canada, Arctic research", "America/Edmonton"},
	"thule":         {76.53, -68.7, "Remote base in Greenland, Arctic monitoring station", "America/Godthab"},
	"barrow":        {71.29, -156.79, "Northernmost city in USA, Arctic coastal town", "America/Anchorage"},
	"longyearbyen":  {78.22, 15.63, "Norwegian settlement in Svalbard, Polar explorer base", "Europe/Oslo"},
	"tromso":        {69.65, 18.96, "Gateway to Arctic, Norwegian polar city", "Europe/Oslo"},
	"marmaris":      {36.8529, 28.2744, "Tropical resort city in Turkey, Mediterranean coast", "Europe/Istanbul"},
}

// getTimezone returns the timezone location for a given IANA timezone string
func getTimezone(tz string) *time.Location {
	if tz == "" || tz == "UTC" {
		return time.UTC
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return time.UTC
	}
	return loc
}

