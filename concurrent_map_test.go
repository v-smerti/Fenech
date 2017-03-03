package fenech

import (
	"encoding/json"
	"hash/fnv"
	"sort"
	"strconv"
	"testing"
)

type Animal []byte

func isAnimal(sl1, sl2 Animal) bool {
	for key, value := range sl1 {
		if sl2[key] != value {
			return false
		}
	}
	return true
}

var m *Fenech

func init() {
	M, err := New("./test")
	if err != nil {
		panic(err)
	}
	m = M

}
func TestMapCreation(t *testing.T) {

	if m == nil {
		t.Error("map is null.")
	}

	if m.Count() != 0 {
		t.Error("new map should be empty.")
	}
}

func TestInsert(t *testing.T) {
	elephant := Animal("elephant")
	monkey := Animal("monkey")

	m.Set("elephant", elephant)
	m.Set("monkey", monkey)

	if m.Count() != 2 {
		t.Error("map should contain exactly two elements.")
	}
}

func TestInsertAbsent(t *testing.T) {
	elephant := Animal("elephant")
	monkey := Animal("monkey")

	m.SetIfAbsent("elephant", elephant)
	ok, err := m.SetIfAbsent("elephant", monkey)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Error("map set a new value even the entry is already present")
	}
}

func TestGet(t *testing.T) {

	// Get a missing element.
	val, ok := m.Get("Money")

	if ok == true {
		t.Error("ok should be false when item is missing from map.")
	}

	if val != nil {
		t.Error("Missing values should return as null.")
	}

	elephant := Animal("elephant")
	m.Set("elephant", elephant)

	// Retrieve inserted element.

	m1, ok := m.Get("elephant")

	if ok == false {
		t.Error("ok should be true for item stored within the map.")
	}

	if &elephant == nil {
		t.Error("expecting an element, not null.")
	}

	if !isAnimal(elephant, m1) {
		t.Log(elephant)
		t.Log(m1)
		t.Error("item was modified.")
	}
}

func TestHas(t *testing.T) {

	// Get a missing element.
	if m.Has("Money") == true {
		t.Error("element shouldn't exists")
	}

	elephant := Animal("elephant")
	m.Set("elephant", elephant)

	if m.Has("elephant") == false {
		t.Error("element exists, expecting Has to return True.")
	}
}

func TestRemove(t *testing.T) {
	removeAllKey()
	monkey := Animal("monkey")
	m.Set("monkey", monkey)

	m.Remove("monkey")

	if m.Count() != 0 {
		t.Error("Expecting count to be zero once item was removed.")
	}

	temp, ok := m.Get("monkey")

	if ok != false {
		t.Error("Expecting ok to be false for missing items.")
	}

	if temp != nil {
		t.Error("Expecting item to be nil after its removal.")
	}

	// Remove a none existing element.
	m.Remove("noone")
}

func TestPop(t *testing.T) {

	monkey := Animal("monkey")
	m.Set("monkey", monkey)

	ml, exists, err := m.Pop("monkey")
	if err != nil {
		t.Fatal(err)
	}

	if !exists {
		t.Error("Pop didn't find a monkey.")
	}

	_, exists2, err := m.Pop("monkey")
	if err != nil {
		t.Fatal(err)
	}
	if exists2 || !isAnimal(ml, monkey) {
		t.Error("Pop keeps finding monkey")
	}

	if m.Count() != 0 {
		t.Error("Expecting count to be zero once item was Pop'ed.")
	}

	temp, ok := m.Get("monkey")

	if ok != false {
		t.Error("Expecting ok to be false for missing items.")
	}

	if temp != nil {
		t.Error("Expecting item to be nil after its removal.")
	}
}

func TestCount(t *testing.T) {
	for i := 0; i < 100; i++ {
		m.Set(strconv.Itoa(i), Animal(strconv.Itoa(i)))
	}

	if m.Count() != 100 {
		t.Error("Expecting 100 element within map.")
	}
}

func TestIsEmpty(t *testing.T) {
	removeAllKey()
	if m.IsEmpty() == false {
		t.Error("new map should be empty")
	}

	m.Set("elephant", Animal("elephant"))

	if m.IsEmpty() != false {
		t.Error("map shouldn't be empty.")
	}
}

func TestIterator(t *testing.T) {
	removeAllKey()
	// Insert 100 elements.
	for i := 0; i < 100; i++ {
		m.Set(strconv.Itoa(i), Animal(strconv.Itoa(i)))
	}

	counter := 0
	// Iterate over elements.
	for item := range m.Iter() {
		val := item.Val

		if val == nil {
			t.Error("Expecting an object.")
		}
		counter++
	}

	if counter != 100 {
		t.Error("We should have counted 100 elements.")
	}
}

func TestBufferedIterator(t *testing.T) {

	// Insert 100 elements.
	for i := 0; i < 100; i++ {
		m.Set(strconv.Itoa(i), Animal(strconv.Itoa(i)))
	}

	counter := 0
	// Iterate over elements.
	for item := range m.IterBuffered() {
		val := item.Val

		if val == nil {
			t.Error("Expecting an object.")
		}
		counter++
	}

	if counter != 100 {
		t.Error("We should have counted 100 elements.")
	}
}

func TestIterCb(t *testing.T) {

	// Insert 100 elements.
	for i := 0; i < 100; i++ {
		m.Set(strconv.Itoa(i), Animal(strconv.Itoa(i)))
	}

	counter := 0
	// Iterate over elements.
	m.IterCb(func(key string, v []byte) {
		counter++
	})
	if counter != 100 {
		t.Error("We should have counted 100 elements.")
	}
}

func TestItems(t *testing.T) {

	// Insert 100 elements.
	for i := 0; i < 100; i++ {
		m.Set(strconv.Itoa(i), Animal(strconv.Itoa(i)))
	}

	items := m.Items()

	if len(items) != 100 {
		t.Error("We should have counted 100 elements.")
	}
}

func TestConcurrent(t *testing.T) {
	ch := make(chan int)
	const iterations = 1000
	var a [iterations]int

	// Using go routines insert 1000 ints into our map.
	go func() {
		for i := 0; i < iterations/2; i++ {
			// Add item to map.
			m.Set(strconv.Itoa(i), Animal(strconv.Itoa(i)))

			// Retrieve item from map.
			val, _ := m.Get(strconv.Itoa(i))
			// Parse int
			k, err := strconv.Atoi(string(val))
			if err != nil {
				t.Error("Err parse int")
			}
			// Write to channel inserted value.
			ch <- k
		} // Call go routine with current index.
	}()

	go func() {
		for i := iterations / 2; i < iterations; i++ {
			// Add item to map.
			m.Set(strconv.Itoa(i), Animal(strconv.Itoa(i)))

			// Retrieve item from map.
			val, _ := m.Get(strconv.Itoa(i))
			// Parse int
			k, err := strconv.Atoi(string(val))
			if err != nil {
				t.Error("Err parse int")
			}
			// Write to channel inserted value.
			ch <- k
		} // Call go routine with current index.
	}()

	// Wait for all go routines to finish.
	counter := 0
	for elem := range ch {
		a[counter] = elem
		counter++
		if counter == iterations {
			break
		}
	}

	// Sorts array, will make is simpler to verify all inserted values we're returned.
	sort.Ints(a[0:iterations])

	// Make sure map contains 1000 elements.
	if m.Count() != iterations {
		t.Error("Expecting 1000 elements.")
	}

	// Make sure all inserted values we're fetched from map.
	for i := 0; i < iterations; i++ {
		if i != a[i] {
			t.Error("missing value", i)
		}
	}
}

func TestJsonMarshal(t *testing.T) {

	expected := `{"0":"MA==","1":"MQ==","10":"MTA=","100":"MTAw","101":"MTAx","102":"MTAy","103":"MTAz","104":"MTA0","105":"MTA1","106":"MTA2","107":"MTA3","108":"MTA4","109":"MTA5","11":"MTE=","110":"MTEw","111":"MTEx","112":"MTEy","113":"MTEz","114":"MTE0","115":"MTE1","116":"MTE2","117":"MTE3","118":"MTE4","119":"MTE5","12":"MTI=","120":"MTIw","121":"MTIx","122":"MTIy","123":"MTIz","124":"MTI0","125":"MTI1","126":"MTI2","127":"MTI3","128":"MTI4","129":"MTI5","13":"MTM=","130":"MTMw","131":"MTMx","132":"MTMy","133":"MTMz","134":"MTM0","135":"MTM1","136":"MTM2","137":"MTM3","138":"MTM4","139":"MTM5","14":"MTQ=","140":"MTQw","141":"MTQx","142":"MTQy","143":"MTQz","144":"MTQ0","145":"MTQ1","146":"MTQ2","147":"MTQ3","148":"MTQ4","149":"MTQ5","15":"MTU=","150":"MTUw","151":"MTUx","152":"MTUy","153":"MTUz","154":"MTU0","155":"MTU1","156":"MTU2","157":"MTU3","158":"MTU4","159":"MTU5","16":"MTY=","160":"MTYw","161":"MTYx","162":"MTYy","163":"MTYz","164":"MTY0","165":"MTY1","166":"MTY2","167":"MTY3","168":"MTY4","169":"MTY5","17":"MTc=","170":"MTcw","171":"MTcx","172":"MTcy","173":"MTcz","174":"MTc0","175":"MTc1","176":"MTc2","177":"MTc3","178":"MTc4","179":"MTc5","18":"MTg=","180":"MTgw","181":"MTgx","182":"MTgy","183":"MTgz","184":"MTg0","185":"MTg1","186":"MTg2","187":"MTg3","188":"MTg4","189":"MTg5","19":"MTk=","190":"MTkw","191":"MTkx","192":"MTky","193":"MTkz","194":"MTk0","195":"MTk1","196":"MTk2","197":"MTk3","198":"MTk4","199":"MTk5","2":"Mg==","20":"MjA=","200":"MjAw","201":"MjAx","202":"MjAy","203":"MjAz","204":"MjA0","205":"MjA1","206":"MjA2","207":"MjA3","208":"MjA4","209":"MjA5","21":"MjE=","210":"MjEw","211":"MjEx","212":"MjEy","213":"MjEz","214":"MjE0","215":"MjE1","216":"MjE2","217":"MjE3","218":"MjE4","219":"MjE5","22":"MjI=","220":"MjIw","221":"MjIx","222":"MjIy","223":"MjIz","224":"MjI0","225":"MjI1","226":"MjI2","227":"MjI3","228":"MjI4","229":"MjI5","23":"MjM=","230":"MjMw","231":"MjMx","232":"MjMy","233":"MjMz","234":"MjM0","235":"MjM1","236":"MjM2","237":"MjM3","238":"MjM4","239":"MjM5","24":"MjQ=","240":"MjQw","241":"MjQx","242":"MjQy","243":"MjQz","244":"MjQ0","245":"MjQ1","246":"MjQ2","247":"MjQ3","248":"MjQ4","249":"MjQ5","25":"MjU=","250":"MjUw","251":"MjUx","252":"MjUy","253":"MjUz","254":"MjU0","255":"MjU1","256":"MjU2","257":"MjU3","258":"MjU4","259":"MjU5","26":"MjY=","260":"MjYw","261":"MjYx","262":"MjYy","263":"MjYz","264":"MjY0","265":"MjY1","266":"MjY2","267":"MjY3","268":"MjY4","269":"MjY5","27":"Mjc=","270":"Mjcw","271":"Mjcx","272":"Mjcy","273":"Mjcz","274":"Mjc0","275":"Mjc1","276":"Mjc2","277":"Mjc3","278":"Mjc4","279":"Mjc5","28":"Mjg=","280":"Mjgw","281":"Mjgx","282":"Mjgy","283":"Mjgz","284":"Mjg0","285":"Mjg1","286":"Mjg2","287":"Mjg3","288":"Mjg4","289":"Mjg5","29":"Mjk=","290":"Mjkw","291":"Mjkx","292":"Mjky","293":"Mjkz","294":"Mjk0","295":"Mjk1","296":"Mjk2","297":"Mjk3","298":"Mjk4","299":"Mjk5","3":"Mw==","30":"MzA=","300":"MzAw","301":"MzAx","302":"MzAy","303":"MzAz","304":"MzA0","305":"MzA1","306":"MzA2","307":"MzA3","308":"MzA4","309":"MzA5","31":"MzE=","310":"MzEw","311":"MzEx","312":"MzEy","313":"MzEz","314":"MzE0","315":"MzE1","316":"MzE2","317":"MzE3","318":"MzE4","319":"MzE5","32":"MzI=","320":"MzIw","321":"MzIx","322":"MzIy","323":"MzIz","324":"MzI0","325":"MzI1","326":"MzI2","327":"MzI3","328":"MzI4","329":"MzI5","33":"MzM=","330":"MzMw","331":"MzMx","332":"MzMy","333":"MzMz","334":"MzM0","335":"MzM1","336":"MzM2","337":"MzM3","338":"MzM4","339":"MzM5","34":"MzQ=","340":"MzQw","341":"MzQx","342":"MzQy","343":"MzQz","344":"MzQ0","345":"MzQ1","346":"MzQ2","347":"MzQ3","348":"MzQ4","349":"MzQ5","35":"MzU=","350":"MzUw","351":"MzUx","352":"MzUy","353":"MzUz","354":"MzU0","355":"MzU1","356":"MzU2","357":"MzU3","358":"MzU4","359":"MzU5","36":"MzY=","360":"MzYw","361":"MzYx","362":"MzYy","363":"MzYz","364":"MzY0","365":"MzY1","366":"MzY2","367":"MzY3","368":"MzY4","369":"MzY5","37":"Mzc=","370":"Mzcw","371":"Mzcx","372":"Mzcy","373":"Mzcz","374":"Mzc0","375":"Mzc1","376":"Mzc2","377":"Mzc3","378":"Mzc4","379":"Mzc5","38":"Mzg=","380":"Mzgw","381":"Mzgx","382":"Mzgy","383":"Mzgz","384":"Mzg0","385":"Mzg1","386":"Mzg2","387":"Mzg3","388":"Mzg4","389":"Mzg5","39":"Mzk=","390":"Mzkw","391":"Mzkx","392":"Mzky","393":"Mzkz","394":"Mzk0","395":"Mzk1","396":"Mzk2","397":"Mzk3","398":"Mzk4","399":"Mzk5","4":"NA==","40":"NDA=","400":"NDAw","401":"NDAx","402":"NDAy","403":"NDAz","404":"NDA0","405":"NDA1","406":"NDA2","407":"NDA3","408":"NDA4","409":"NDA5","41":"NDE=","410":"NDEw","411":"NDEx","412":"NDEy","413":"NDEz","414":"NDE0","415":"NDE1","416":"NDE2","417":"NDE3","418":"NDE4","419":"NDE5","42":"NDI=","420":"NDIw","421":"NDIx","422":"NDIy","423":"NDIz","424":"NDI0","425":"NDI1","426":"NDI2","427":"NDI3","428":"NDI4","429":"NDI5","43":"NDM=","430":"NDMw","431":"NDMx","432":"NDMy","433":"NDMz","434":"NDM0","435":"NDM1","436":"NDM2","437":"NDM3","438":"NDM4","439":"NDM5","44":"NDQ=","440":"NDQw","441":"NDQx","442":"NDQy","443":"NDQz","444":"NDQ0","445":"NDQ1","446":"NDQ2","447":"NDQ3","448":"NDQ4","449":"NDQ5","45":"NDU=","450":"NDUw","451":"NDUx","452":"NDUy","453":"NDUz","454":"NDU0","455":"NDU1","456":"NDU2","457":"NDU3","458":"NDU4","459":"NDU5","46":"NDY=","460":"NDYw","461":"NDYx","462":"NDYy","463":"NDYz","464":"NDY0","465":"NDY1","466":"NDY2","467":"NDY3","468":"NDY4","469":"NDY5","47":"NDc=","470":"NDcw","471":"NDcx","472":"NDcy","473":"NDcz","474":"NDc0","475":"NDc1","476":"NDc2","477":"NDc3","478":"NDc4","479":"NDc5","48":"NDg=","480":"NDgw","481":"NDgx","482":"NDgy","483":"NDgz","484":"NDg0","485":"NDg1","486":"NDg2","487":"NDg3","488":"NDg4","489":"NDg5","49":"NDk=","490":"NDkw","491":"NDkx","492":"NDky","493":"NDkz","494":"NDk0","495":"NDk1","496":"NDk2","497":"NDk3","498":"NDk4","499":"NDk5","5":"NQ==","50":"NTA=","500":"NTAw","501":"NTAx","502":"NTAy","503":"NTAz","504":"NTA0","505":"NTA1","506":"NTA2","507":"NTA3","508":"NTA4","509":"NTA5","51":"NTE=","510":"NTEw","511":"NTEx","512":"NTEy","513":"NTEz","514":"NTE0","515":"NTE1","516":"NTE2","517":"NTE3","518":"NTE4","519":"NTE5","52":"NTI=","520":"NTIw","521":"NTIx","522":"NTIy","523":"NTIz","524":"NTI0","525":"NTI1","526":"NTI2","527":"NTI3","528":"NTI4","529":"NTI5","53":"NTM=","530":"NTMw","531":"NTMx","532":"NTMy","533":"NTMz","534":"NTM0","535":"NTM1","536":"NTM2","537":"NTM3","538":"NTM4","539":"NTM5","54":"NTQ=","540":"NTQw","541":"NTQx","542":"NTQy","543":"NTQz","544":"NTQ0","545":"NTQ1","546":"NTQ2","547":"NTQ3","548":"NTQ4","549":"NTQ5","55":"NTU=","550":"NTUw","551":"NTUx","552":"NTUy","553":"NTUz","554":"NTU0","555":"NTU1","556":"NTU2","557":"NTU3","558":"NTU4","559":"NTU5","56":"NTY=","560":"NTYw","561":"NTYx","562":"NTYy","563":"NTYz","564":"NTY0","565":"NTY1","566":"NTY2","567":"NTY3","568":"NTY4","569":"NTY5","57":"NTc=","570":"NTcw","571":"NTcx","572":"NTcy","573":"NTcz","574":"NTc0","575":"NTc1","576":"NTc2","577":"NTc3","578":"NTc4","579":"NTc5","58":"NTg=","580":"NTgw","581":"NTgx","582":"NTgy","583":"NTgz","584":"NTg0","585":"NTg1","586":"NTg2","587":"NTg3","588":"NTg4","589":"NTg5","59":"NTk=","590":"NTkw","591":"NTkx","592":"NTky","593":"NTkz","594":"NTk0","595":"NTk1","596":"NTk2","597":"NTk3","598":"NTk4","599":"NTk5","6":"Ng==","60":"NjA=","600":"NjAw","601":"NjAx","602":"NjAy","603":"NjAz","604":"NjA0","605":"NjA1","606":"NjA2","607":"NjA3","608":"NjA4","609":"NjA5","61":"NjE=","610":"NjEw","611":"NjEx","612":"NjEy","613":"NjEz","614":"NjE0","615":"NjE1","616":"NjE2","617":"NjE3","618":"NjE4","619":"NjE5","62":"NjI=","620":"NjIw","621":"NjIx","622":"NjIy","623":"NjIz","624":"NjI0","625":"NjI1","626":"NjI2","627":"NjI3","628":"NjI4","629":"NjI5","63":"NjM=","630":"NjMw","631":"NjMx","632":"NjMy","633":"NjMz","634":"NjM0","635":"NjM1","636":"NjM2","637":"NjM3","638":"NjM4","639":"NjM5","64":"NjQ=","640":"NjQw","641":"NjQx","642":"NjQy","643":"NjQz","644":"NjQ0","645":"NjQ1","646":"NjQ2","647":"NjQ3","648":"NjQ4","649":"NjQ5","65":"NjU=","650":"NjUw","651":"NjUx","652":"NjUy","653":"NjUz","654":"NjU0","655":"NjU1","656":"NjU2","657":"NjU3","658":"NjU4","659":"NjU5","66":"NjY=","660":"NjYw","661":"NjYx","662":"NjYy","663":"NjYz","664":"NjY0","665":"NjY1","666":"NjY2","667":"NjY3","668":"NjY4","669":"NjY5","67":"Njc=","670":"Njcw","671":"Njcx","672":"Njcy","673":"Njcz","674":"Njc0","675":"Njc1","676":"Njc2","677":"Njc3","678":"Njc4","679":"Njc5","68":"Njg=","680":"Njgw","681":"Njgx","682":"Njgy","683":"Njgz","684":"Njg0","685":"Njg1","686":"Njg2","687":"Njg3","688":"Njg4","689":"Njg5","69":"Njk=","690":"Njkw","691":"Njkx","692":"Njky","693":"Njkz","694":"Njk0","695":"Njk1","696":"Njk2","697":"Njk3","698":"Njk4","699":"Njk5","7":"Nw==","70":"NzA=","700":"NzAw","701":"NzAx","702":"NzAy","703":"NzAz","704":"NzA0","705":"NzA1","706":"NzA2","707":"NzA3","708":"NzA4","709":"NzA5","71":"NzE=","710":"NzEw","711":"NzEx","712":"NzEy","713":"NzEz","714":"NzE0","715":"NzE1","716":"NzE2","717":"NzE3","718":"NzE4","719":"NzE5","72":"NzI=","720":"NzIw","721":"NzIx","722":"NzIy","723":"NzIz","724":"NzI0","725":"NzI1","726":"NzI2","727":"NzI3","728":"NzI4","729":"NzI5","73":"NzM=","730":"NzMw","731":"NzMx","732":"NzMy","733":"NzMz","734":"NzM0","735":"NzM1","736":"NzM2","737":"NzM3","738":"NzM4","739":"NzM5","74":"NzQ=","740":"NzQw","741":"NzQx","742":"NzQy","743":"NzQz","744":"NzQ0","745":"NzQ1","746":"NzQ2","747":"NzQ3","748":"NzQ4","749":"NzQ5","75":"NzU=","750":"NzUw","751":"NzUx","752":"NzUy","753":"NzUz","754":"NzU0","755":"NzU1","756":"NzU2","757":"NzU3","758":"NzU4","759":"NzU5","76":"NzY=","760":"NzYw","761":"NzYx","762":"NzYy","763":"NzYz","764":"NzY0","765":"NzY1","766":"NzY2","767":"NzY3","768":"NzY4","769":"NzY5","77":"Nzc=","770":"Nzcw","771":"Nzcx","772":"Nzcy","773":"Nzcz","774":"Nzc0","775":"Nzc1","776":"Nzc2","777":"Nzc3","778":"Nzc4","779":"Nzc5","78":"Nzg=","780":"Nzgw","781":"Nzgx","782":"Nzgy","783":"Nzgz","784":"Nzg0","785":"Nzg1","786":"Nzg2","787":"Nzg3","788":"Nzg4","789":"Nzg5","79":"Nzk=","790":"Nzkw","791":"Nzkx","792":"Nzky","793":"Nzkz","794":"Nzk0","795":"Nzk1","796":"Nzk2","797":"Nzk3","798":"Nzk4","799":"Nzk5","8":"OA==","80":"ODA=","800":"ODAw","801":"ODAx","802":"ODAy","803":"ODAz","804":"ODA0","805":"ODA1","806":"ODA2","807":"ODA3","808":"ODA4","809":"ODA5","81":"ODE=","810":"ODEw","811":"ODEx","812":"ODEy","813":"ODEz","814":"ODE0","815":"ODE1","816":"ODE2","817":"ODE3","818":"ODE4","819":"ODE5","82":"ODI=","820":"ODIw","821":"ODIx","822":"ODIy","823":"ODIz","824":"ODI0","825":"ODI1","826":"ODI2","827":"ODI3","828":"ODI4","829":"ODI5","83":"ODM=","830":"ODMw","831":"ODMx","832":"ODMy","833":"ODMz","834":"ODM0","835":"ODM1","836":"ODM2","837":"ODM3","838":"ODM4","839":"ODM5","84":"ODQ=","840":"ODQw","841":"ODQx","842":"ODQy","843":"ODQz","844":"ODQ0","845":"ODQ1","846":"ODQ2","847":"ODQ3","848":"ODQ4","849":"ODQ5","85":"ODU=","850":"ODUw","851":"ODUx","852":"ODUy","853":"ODUz","854":"ODU0","855":"ODU1","856":"ODU2","857":"ODU3","858":"ODU4","859":"ODU5","86":"ODY=","860":"ODYw","861":"ODYx","862":"ODYy","863":"ODYz","864":"ODY0","865":"ODY1","866":"ODY2","867":"ODY3","868":"ODY4","869":"ODY5","87":"ODc=","870":"ODcw","871":"ODcx","872":"ODcy","873":"ODcz","874":"ODc0","875":"ODc1","876":"ODc2","877":"ODc3","878":"ODc4","879":"ODc5","88":"ODg=","880":"ODgw","881":"ODgx","882":"ODgy","883":"ODgz","884":"ODg0","885":"ODg1","886":"ODg2","887":"ODg3","888":"ODg4","889":"ODg5","89":"ODk=","890":"ODkw","891":"ODkx","892":"ODky","893":"ODkz","894":"ODk0","895":"ODk1","896":"ODk2","897":"ODk3","898":"ODk4","899":"ODk5","9":"OQ==","90":"OTA=","900":"OTAw","901":"OTAx","902":"OTAy","903":"OTAz","904":"OTA0","905":"OTA1","906":"OTA2","907":"OTA3","908":"OTA4","909":"OTA5","91":"OTE=","910":"OTEw","911":"OTEx","912":"OTEy","913":"OTEz","914":"OTE0","915":"OTE1","916":"OTE2","917":"OTE3","918":"OTE4","919":"OTE5","92":"OTI=","920":"OTIw","921":"OTIx","922":"OTIy","923":"OTIz","924":"OTI0","925":"OTI1","926":"OTI2","927":"OTI3","928":"OTI4","929":"OTI5","93":"OTM=","930":"OTMw","931":"OTMx","932":"OTMy","933":"OTMz","934":"OTM0","935":"OTM1","936":"OTM2","937":"OTM3","938":"OTM4","939":"OTM5","94":"OTQ=","940":"OTQw","941":"OTQx","942":"OTQy","943":"OTQz","944":"OTQ0","945":"OTQ1","946":"OTQ2","947":"OTQ3","948":"OTQ4","949":"OTQ5","95":"OTU=","950":"OTUw","951":"OTUx","952":"OTUy","953":"OTUz","954":"OTU0","955":"OTU1","956":"OTU2","957":"OTU3","958":"OTU4","959":"OTU5","96":"OTY=","960":"OTYw","961":"OTYx","962":"OTYy","963":"OTYz","964":"OTY0","965":"OTY1","966":"OTY2","967":"OTY3","968":"OTY4","969":"OTY5","97":"OTc=","970":"OTcw","971":"OTcx","972":"OTcy","973":"OTcz","974":"OTc0","975":"OTc1","976":"OTc2","977":"OTc3","978":"OTc4","979":"OTc5","98":"OTg=","980":"OTgw","981":"OTgx","982":"OTgy","983":"OTgz","984":"OTg0","985":"OTg1","986":"OTg2","987":"OTg3","988":"OTg4","989":"OTg5","99":"OTk=","990":"OTkw","991":"OTkx","992":"OTky","993":"OTkz","994":"OTk0","995":"OTk1","996":"OTk2","997":"OTk3","998":"OTk4","999":"OTk5","a":"MQ==","b":"Mg=="}`

	m.Set("a", []byte("1"))
	m.Set("b", []byte("2"))
	j, err := json.Marshal(m)
	if err != nil {
		t.Error(err)
	}

	if string(j) != expected {
		t.Error("json", string(j), "differ from expected", expected)
		return
	}
}

func TestKeys(t *testing.T) {
	removeAllKey()
	// Insert 100 elements.
	for i := 0; i < 100; i++ {
		m.Set(strconv.Itoa(i), Animal(strconv.Itoa(i)))
	}

	keys := m.Keys()
	if len(keys) != 100 {
		t.Error("We should have counted 100 elements.")
	}
}

func TestMInsert(t *testing.T) {
	removeAllKey()
	animals := map[string][]byte{
		"elephant": Animal("elephant"),
		"monkey":   Animal("monkey"),
	}

	m.MSet(animals)

	if m.Count() != 2 {
		t.Error("map should contain exactly two elements.")
	}
}

func TestFnv32(t *testing.T) {
	key := []byte("ABC")

	hasher := fnv.New32()
	hasher.Write(key)
	if fnv32(string(key)) != hasher.Sum32() {
		t.Errorf("Bundled fnv32 produced %d, expected result from hash/fnv32 is %d", fnv32(string(key)), hasher.Sum32())
	}
}

func TestUpsert(t *testing.T) {
	removeAllKey()
	dolphin := Animal("dolphin")
	whale := Animal("whale")
	tiger := Animal("tiger")
	lion := Animal("lion")

	cb := func(exists bool, valueInMap []byte, newValue []byte) []byte {
		if !exists {
			return []byte("cb_")
		}
		return append(valueInMap, newValue...)
	}

	m.Set("marine", []byte(dolphin))
	m.Upsert("marine", whale, cb)
	m.Upsert("predator", tiger, cb)
	m.Upsert("predator", lion, cb)

	if m.Count() != 2 {
		t.Error("map should contain exactly two elements.")
	}

	marineAnimals, ok := m.Get("marine")

	if !ok || string(marineAnimals) != "dolphinwhale" {
		t.Error("Set, then Upsert failed")
	}

	predators, ok := m.Get("predator")
	if !ok || string(predators) != "cb_lion" {
		t.Error("Upsert, then Upsert failed")
	}

}

func TestKeysWhenRemoving(t *testing.T) {

	// Insert 100 elements.
	Total := 100
	for i := 0; i < Total; i++ {
		m.Set(strconv.Itoa(i), Animal(strconv.Itoa(i)))
	}

	// Remove 10 elements concurrently.
	Num := 10
	for i := 0; i < Num; i++ {
		go func(c *Fenech, n int) {
			c.Remove(strconv.Itoa(n))
		}(m, i)
	}
	keys := m.Keys()
	for _, k := range keys {
		if k == "" {
			t.Error("Empty keys returned")
		}
	}
}

func removeAllKey() {
	for _, key := range m.Keys() {
		m.Remove(key)
	}
}
