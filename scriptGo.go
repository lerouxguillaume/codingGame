package main

import "fmt"
//import "os"
import "sort"
//import "math/rand"

/**
 * Auto-generated code below aims at helping you parse
 * the standard input according to the problem statement.
 **/
const MAX_PREDICT = 15;
const INFINITE = 1000;
const SAFE_ZONE = 3;
const RELOAD_BOMB = 20;
var reload = 0;
var turn = 0;

// AxisSorter sorts planets by axis.
type RouteSorter []Route

func (a RouteSorter) Len() int           { return len(a) }
func (a RouteSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a RouteSorter) Less(i, j int) bool { return a[i].cost < a[j].cost }

type Spare struct {
        unit int,
             delay int
}

type Factory struct {
        id int,
           camp int,
           units [MAX_PREDICT]int,
           prod int,
           disabled int,
           spare Spare,
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
        return fmt.Sprint("MOVE ",f.id,f.shortestPath(target),troups,";");
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

func (f *Factory)makeFactory(id int,camp int, units int, prod int,disabled int) {
        f.id = id;
        f.camp = camp;
        f.prod = prod;
        f.disabled = disabled;
        if(camp == 0){prod = 0;} //Neutral camp doesn't produce
        for key, _ := range res.units {
                if(key > disabled){
                        f.units[key] = units+(key-disabled)*prod;
                }else{
                        f.units[key] = units
                }
        };
}

func makeArmy(camp int, start int, target int,units int, delay int) Army {
        return Army{camp:camp,start:start,target:target,units:units,delay:delay}
}

type Army struct {
        camp int,
             start int,
             target int,
             units int,
             delay int
}

type Route struct{
        start *Factory,
              end *Factory,
              cost int
}


type Graph struct{
        nbVertex int,
                 nbRoute int,
                 routes []Route,
                 factorys []Factory,
                 armys []Army,
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

func (f *Factory) evaluateDistToEnenmys()(int){
count := 0;
       for _,value := range f.routes{
               if(value.end.camp == -1){
                       count +=  value.cost;
               }
       }
       return count
}

func (f *Factory) sendToFront()(string){
frontLane := Target{id:-1}
best := INFINITE;
      for _,value := range f.routes{
tmp := value.end.evaluateDistToEnenmys()
             if(tmp < best){
                     best = tmp;
                     frontLane = Target{id:value.end.id,units:f.spare.unit}
             }
      }
      //fmt.Fprintln(os.Stderr, "debug",frontLane)
      return f.sendUnits(frontLane.id,frontLane.units);
}

func (g *Graph) expand(){
        for i:=0;i<len(g.factorys);i++ {
                if(g.factorys[i].camp == 1 && g.factorys[i].spare.unit>0){
current := g.factorys[i]
                 target := g.lookForTargets(current)
                 for k := 0; k < len(target); k++ {
                         //fmt.Fprintln(os.Stderr, "need help of ",value.spare.unit,"i ask : ", -value.spare.unit)		
                         g.msg += current.sendUnits(target[k].id, target[k].units);
                 }
         if(current.spare.unit > 30 && current.prod < 3)||(current.prod == 0 && current.spare.unit >= 10){
                 g.msg += fmt.Sprint("INC ",current.id,";")
         }else if(current.prod == 3){
                 g.msg += g.factorys[i].sendToFront()
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
                          reload = RELOAD_BOMB;
          }
}


func (f *Factory)amICloser()bool{
they := INFINITE;
me := INFINITE;
    for _,value := range f.routes {
            if(value.end.camp == 1 && value.cost < me){
                    me = value.cost
            }else if(value.end.camp == -1 && value.cost < they){
                    they = value.cost;
            }
    }
    return me < they;
}

func (g *Graph)spread(){
        var targets []Factory;
        var hunter Factory;
        for i:=0;i<len(g.factorys);i++ {
                if(g.factorys[i].camp == 1){
                        hunter = g.factorys[i];
                        for _,value := range hunter.routes{
                                if(value.end.camp == 0 && value.end.prod > 0 &&  value.end.amICloser()){
                                        targets= append(targets, *value.end);
                                }
                        }
                }
        }
soldiers := hunter.spare.unit;
          targets = sacADos(soldiers,targets);
          for _,target := range targets{
                  g.msg += hunter.sendUnits(target.id,target.units[0]+1);
          }
}

func sacADos(capacity int, objects []Factory) []Factory{
        var res []Factory
                if(len(objects) == 1){
                        return append(res,objects[0]);
                }else if(len(objects) < 1){
                        return res
                }else{
m := make([][]int,len(objects));
   for i := range m{
           m[i] = make([]int, capacity);
   }
   for j:= 0;j < capacity;j++{
           if(objects[0].units[0]+1 > j){
                   m[0][j] = 0
           }else{
                   m[0][j] = objects[0].prod;
           }
   }
   for i := 1; i < len(objects);i++{
           for j:= 0;j < capacity;j++{
                   if(objects[i].units[0]+1 > j){
                           m[i][j] = m[i-1][j]
                   }else{
arg1 := m[i-1][j];
arg2 := m[i-1][j-objects[i].units[0]-1] + objects[i].prod;
      if(arg1 > arg2){
              m[i][j] = arg1
      }else{
              m[i][j] = arg2
      }
                   }
           }
   }

j:= capacity-1
          i:= len(objects)-1;
  //fmt.Fprintln(os.Stderr, "DEBUG START ", j);

  for (m[i][j] == m[i][j-1]){
          j--
  }
  //fmt.Fprintln(os.Stderr, "DEBUG START ", j);

  for(j > 0){
          for(i > 0 && m[i][j] == m[i-1][j]){
                  i--
          }
          j -= objects[i].units[0]+1;
          //fmt.Fprintln(os.Stderr, "DEBUG", j);

          if(j >= 0){
                  res = append(res, objects[i]);
          }
          i--
  }
                }
        //fmt.Fprintln(os.Stderr, "objects", len(objects), objects);
        return res;
}

func (g *Graph) play() {
        g.state();
        if(turn > 0){
                g.help();
                g.expand();
        }else{
                g.spread();
        }
        g.bomb();
        g.msg += "WAIT"
                turn++;
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
                                route := Route{start:&graph.factorys[factory1],end:&graph.factorys[factory2],cost:distance}
                        graph.routes[i] = route;
                        route.end.routes = append(route.end.routes, Route{start:route.end,end:route.start,cost:route.cost})
                                route.start.routes = append(route.start.routes, route)
                }
        for _,fact := range graph.factorys{
                sort.Sort(RouteSorter(fact.routes))
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
                                        graph.factorys[i].makeFactory(entityId,arg1,arg2,arg3,arg4);
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
