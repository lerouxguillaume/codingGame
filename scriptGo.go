		package main

		import "fmt"
		//import "os"
		import "sort"
		import "math/rand"

		/**
		 * Auto-generated code below aims at helping you parse
		 * the standard input according to the problem statement.
		 **/
		const MAX_PREDICT = 15;
		const INFINITE = 1000;
		const SAFE_ZONE = 3;
		var reload = 0;

		// AxisSorter sorts planets by axis.
		type RouteSorter []Route

		func (a RouteSorter) Len() int           { return len(a) }
		func (a RouteSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
		func (a RouteSorter) Less(i, j int) bool { return a[i].cost < a[j].cost }
			
		type Spare struct {
			unit int
			delay int
		}

		type Factory struct {
			id int
			camp int
			units [MAX_PREDICT]int
			prod int
			disabled int
			spare Spare
			routes []Route
		}

		func (f *Factory)updateUnits(nbOpponents int, delay int){
			for i:=delay;i<MAX_PREDICT;i++ {
				f.units[i] -= nbOpponents
			}
		}
		
		func (f *Factory)setRoutes(routes []Route){
			for _,value := range routes {
				if(value.end.id == f.id){
					f.routes = append(f.routes, Route{start:value.end,end:value.start,cost:value.cost})
				}
				if(value.start.id == f.id){
					f.routes = append(f.routes, value)
				}
			}
			sort.Sort(RouteSorter(f.routes))
		}
		
		func (f *Factory)sendUnits(target int, units int) string {
			troups := units
			if(f.spare.unit < 0){
				troups = 0;
			} else if(units > f.spare.unit){
				troups = f.spare.unit;
			}
			f.spare.unit -= troups;
			//fmt.Fprintln(os.Stderr, "need help",units,"send : ", troups,"Available (spare)", f.spare.unit)		
			return fmt.Sprint("MOVE ",f.id,f.shortestPath(target),troups,";") 
		}

		func (f *Factory)shortestPath(target int)(int){
			short := target;
			dist := 21;
			for _,value := range f.routes{
				if(value.end.id == target){
					dist = value.cost;
				}
			}
			for _,value := range f.routes {
				//fmt.Fprintln(os.Stderr, "I'M : ",f.id, " TARGETING : ", value.end.id, "COST : ",value.cost)		
				tmp := recursive(*value.end, target,dist,value.cost)
				if(tmp < dist){
					dist = tmp;
					short = value.end.id;
				}
			}
			return short;
		}
		
		func recursive(f Factory, target int, maxCost int, currentCost int) (int){
			min := maxCost;
			if(f.id == target){
				return currentCost;
			}else{
				//currentCost++;
				for _,value := range f.routes {
					if(currentCost+value.cost < maxCost){
						tmp := recursive(*value.end, target, maxCost,currentCost+value.cost);
						if(tmp <= min){
							min = tmp;
						}
					}
				}
			}
			return min;
		}
		
		func makeFactory(id int,camp int, units int, prod int,disabled int) Factory {
			var res = Factory{id:id,camp:camp,prod:prod,disabled:disabled}
			if(camp == 0){prod = 0;} //Neutral camp doesn't produce
			for key, _ := range res.units {
				if(key > disabled){
					res.units[key] = units+(key-disabled)*prod;
				}else{
					res.units[key] = units
				}
			}	
			return res;
		}

		func makeArmy(camp int, start int, target int,units int, delay int) Army {
			return Army{camp:camp,start:start,target:target,units:units,delay:delay}
		}

		type Army struct {
			camp int
			start int
			target int
			units int
			delay int
		}

		type Route struct{
			start *Factory
			end *Factory
			cost int	
		}


		type Graph struct{
			nbVertex int
			nbRoute int
			routes []Route
			factorys []Factory
			armys []Army
			msg string
		}

		func (g *Graph) getFactory(id int) *Factory {
			for i:=0;i<len(g.factorys);i++ {
				if(g.factorys[i].id == id){
					return &g.factorys[i];
				}
			}
			return nil;
		}

		func (g *Graph) getRoutes() []Route {
			return g.routes
		}
		func (g *Graph) getFactorys() []Factory {
			return g.factorys
		}

		func (g *Graph) state(){
			for i:=0;i<len(g.armys);i++ {
				current := g.armys[i];
				if(current.delay < MAX_PREDICT){
					fact := g.getFactory(current.target);
					if(fact.camp == current.camp){
						fact.updateUnits(-current.units,current.delay);
					} else{
						fact.updateUnits(current.units,current.delay);					
					}
				}
			}
			g.calculateSpare();
		}
		
		func (g *Graph) calculateSpare(){
			for i:=0;i<len(g.factorys);i++ {
				current := &g.factorys[i];
				spare := Spare{unit:INFINITE,delay:INFINITE}
				if(current.camp == 1){
					for key,value := range current.units {
						if(value < spare.unit){
							spare.unit = value;
							if(value < 0 && spare.delay == INFINITE){
								spare.delay = key;
							}
						}			    
					}
					current.spare = spare;
				}
				if(g.factorys[i].id == 3){
				//fmt.Fprintln(os.Stderr, "debug EXPAND id: ",g.factorys[i].spare," units : ",g.factorys[i].units);
				}
			}
		}
		
		func (g *Graph)getAllysFrom(id int,dist int) []Factory{
			var res []Factory;
			for _,value := range g.routes {
				if(value.cost == dist){
					if (value.start.id == id && value.end.camp == 1 && value.end.spare.unit > 0){
						res = append(res, *value.end);
					} 
					if (value.end.id == id && value.start.camp == 1 && value.start.spare.unit > 0){
						res = append(res, *value.start);
					}
				}
			}
			return res;
		}
		
		
		func (g *Graph) help(){
			for _,value := range g.factorys {
				helpers := [][]Factory{}
				if(value.camp == 1 && value.spare.unit < 0){
					for i:=1 ; i < value.spare.delay; i++ {
						helpers = append(helpers,g.getAllysFrom(value.id,i))
					}
					for _,helpersDist := range helpers{
						for _,helperDist := range helpersDist{
							helperFact := g.getFactory(helperDist.id);
							//fmt.Fprintln(os.Stderr, "need help of ",value.spare.unit,"i ask : ", -value.spare.unit)		
							g.msg += helperFact.sendUnits(value.id, -value.spare.unit);
						}
					}
				}
			}
		}
		
		type Target struct{
			id int
			units int
		}
		
		
		func (g *Graph) lookForTargets(hunter Factory)([]Target){
			soldiers := hunter.spare.unit
			var a []Target
			safety := 0;
			for _,value := range hunter.routes{
				if(value.end.camp == 1){
					safety++;
					if(safety == SAFE_ZONE && hunter.prod < 3 && hunter.spare.unit > 10){
						g.msg += fmt.Sprint("INC ",hunter.id,";") 
						break;
					}
				}else if(value.end.camp == 0 && value.end.prod > 0 && value.cost < MAX_PREDICT-1){
					safety = 0;
					if(soldiers > value.end.units[value.cost+1] && value.end.units[value.cost+1] >= 0){
						a = append(a,Target{id:value.end.id,units:(value.end.units[0]+1)})
						soldiers -= value.end.units[0]+1
					}
				}else if(value.end.camp == -1 && value.end.prod > 0 && value.cost < MAX_PREDICT-1){
					safety = 0;
					//fmt.Fprintln(os.Stderr, "debug id: ",hunter.id," soldiers : ",soldiers," def : ",value.end.units[value.cost])
					if(soldiers > value.end.units[value.cost+1] && value.end.units[value.cost+1] >= 0){
						a = append(a,Target{id:value.end.id,units:(value.end.units[value.cost+1]+1)})
						soldiers -= value.end.units[value.cost]+1
					}
				}
				//fmt.Fprintln(os.Stderr, "debug",key,value.cost)
			}
			return a
			
		}
		
		func (g *Graph) expand(){
			for i:=0;i<len(g.factorys);i++ {
			if(g.factorys[i].id == 3){
				//fmt.Fprintln(os.Stderr, "debug EXPAND id: ",g.factorys[i].spare);
			}
				if(g.factorys[i].camp == 1 && g.factorys[i].spare.unit>0){
					current := g.factorys[i]
					if(current.spare.unit > 30 && current.prod < 3)||(current.prod == 0 && current.spare.unit >= 10){
						g.msg += fmt.Sprint("INC ",current.id,";") 
					}else if(current.spare.unit > 100 && current.prod == 3){
						random := rand.Intn(len(g.factorys));
						for ( random == current.id){
							random = rand.Intn(len(g.factorys));
						}
						g.msg += current.sendUnits(random, 10);
					}
						target := g.lookForTargets(current)
						for k := 0; k < len(target); k++ {
							//fmt.Fprintln(os.Stderr, "need help of ",value.spare.unit,"i ask : ", -value.spare.unit)		
							g.msg += current.sendUnits(target[k].id, target[k].units);
						}
					
				}
			}
		}
		
		
		func (g *Graph)bomb(){
			reload--
			myProd := 0;
			theirProd := 0;
			minDist := Target{units:INFINITE};
			sourceId := -1;
			for  _,value := range g.factorys{
				if(value.camp == 1){
					myProd += value.prod
					for i := 0; i< len(value.routes); i++ {
						currentRoute := value.routes[i];
						if(currentRoute.end.camp == -1 && minDist.units > currentRoute.cost && currentRoute.end.prod == 3){
							minDist = Target{id:currentRoute.end.id,units : currentRoute.cost}
							sourceId = value.id
						}
					}
				}else if(value.camp == -1){
					theirProd += value.prod;
				}
			}
			if(theirProd >= myProd && reload < 0 && sourceId != -1){
				g.msg += fmt.Sprint("BOMB ",sourceId,minDist.id,";")
				reload = 20;
			}
		}
		
		func (g *Graph) play() {
			g.state();
			g.help();
			g.expand();
			g.bomb();
			g.msg += "WAIT"
		}
		 
		 
		func main() {
			// factoryCount: the number of factories
			var factoryCount int
			fmt.Scan(&factoryCount)
			// linkCount: the number of links between factories
			var linkCount int
			fmt.Scan(&linkCount)
			
			graph := &Graph{}
			graph.routes = make([]Route, linkCount)
			graph.factorys = make([]Factory, factoryCount)

			for i := 0; i < linkCount; i++ {
				var factory1, factory2, distance int
				fmt.Scan(&factory1, &factory2, &distance)
				graph.routes[i] = Route{start:&graph.factorys[factory1],end:&graph.factorys[factory2],cost:distance}
			}
			for {
				// entityCount: the number of entities (e.g. factories and troops)
				var entityCount int
				fmt.Scan(&entityCount)
				graph.armys = nil;

				for i := 0; i < entityCount; i++ {
					var entityId int
					var entityType string
					var arg1, arg2, arg3, arg4, arg5 int
					fmt.Scan(&entityId, &entityType, &arg1, &arg2, &arg3, &arg4, &arg5)
					if(entityType == "FACTORY"){
						graph.factorys[i] = makeFactory(entityId,arg1,arg2,arg3,arg4);
						graph.factorys[i].setRoutes(graph.routes)
					}
					if(entityType == "TROOP"){
						graph.armys = append(graph.armys,makeArmy(arg1,arg2,arg3,arg4,arg5))
					}
				}	
				graph.msg = "";
				graph.play();
				fmt.Println(graph.msg)
				// Any valid action, such as "WAIT" or "MOVE source destination cyborgs"
			}
		}