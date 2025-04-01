export namespace main {
	
	export class ShopItem {
	    index: number;
	    itemId: number;
	    itemName: string;
	    itemCount: number;
	    itemPrice: number;
	    itemDiscount: number;
	    itemDisplayColor: string;
	    itemSpecialType: number;
	
	    static createFrom(source: any = {}) {
	        return new ShopItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.index = source["index"];
	        this.itemId = source["itemId"];
	        this.itemName = source["itemName"];
	        this.itemCount = source["itemCount"];
	        this.itemPrice = source["itemPrice"];
	        this.itemDiscount = source["itemDiscount"];
	        this.itemDisplayColor = source["itemDisplayColor"];
	        this.itemSpecialType = source["itemSpecialType"];
	    }
	}

}

