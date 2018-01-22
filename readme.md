Intuit Quickbooks is a multi-tenant product and serves various regions across geographies. The product is also highly configurable, the configurations controlled via a myriad of properties.  And the properties being highly dynamic in nature, the system needs to be aware of the changes in properties and should react accordingly.  In this context, 

 Design and implement a basic config management system that quickbooks can use
1.	For loading all startup properties (static ones)
2.	For loading feature toggles (dynamic ones)
3.	Change in a common property should be visible in all tenants

Features
1.	There should be a single interface to the Config management system
2.	Ability to provide a default value for a property

Examples
   ConfgManager.getString(“property”) ;
   ConfgManager.getString(“property”,”default value”) ;
   ConfgManager.getBoolean(“property”,false) ;



-- Notes:

To send message to client

{"key":"mani1","value":"mani","mType":2}

// client asking configuration
{"clientID":"mani1","key":"abc","mType":1}


update client

http://localhost:8091/notify/mani/kannamma