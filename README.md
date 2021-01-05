# PhysicsEditor-Box2d.go

PhysicsEditor-Box2d.go is simple library to parse [PhysicsEditor's](https://www.codeandweb.com/physicseditor) default Box2D XML format

It creates [Box2D.go](https://github.com/ByteArena/box2d) structs in following format:

        type ParsedBody struct {
        	Name string
        	Shapes []box2d.B2FixtureDef
        }


## TODO
* anchor point (currently working for 0.5, 0.5)
* masks
* probably more but I don't need more features at the moment