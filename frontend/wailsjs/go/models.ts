export namespace core {
	
	export class AddScanTaskParams {
	    path: string;
	
	    static createFrom(source: any = {}) {
	        return new AddScanTaskParams(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	    }
	}
	export class DeleteScanTasksParams {
	    ids: number[];
	
	    static createFrom(source: any = {}) {
	        return new DeleteScanTasksParams(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ids = source["ids"];
	    }
	}
	export class ExportParams {
	    limit: number;
	
	    static createFrom(source: any = {}) {
	        return new ExportParams(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.limit = source["limit"];
	    }
	}
	export class IdsForm {
	    ids: number[];
	
	    static createFrom(source: any = {}) {
	        return new IdsForm(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ids = source["ids"];
	    }
	}
	export class IdsStatus {
	    ids: number[];
	    status: number;
	
	    static createFrom(source: any = {}) {
	        return new IdsStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ids = source["ids"];
	        this.status = source["status"];
	    }
	}
	export class ResData {
	    status: boolean;
	    msg: string;
	    data: any;
	
	    static createFrom(source: any = {}) {
	        return new ResData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.msg = source["msg"];
	        this.data = source["data"];
	    }
	}
	export class ScanListParams {
	    page: number;
	
	    static createFrom(source: any = {}) {
	        return new ScanListParams(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.page = source["page"];
	    }
	}
	export class SettingData {
	    base_url: string;
	    token: string;
	    http_proxy: string;
	
	    static createFrom(source: any = {}) {
	        return new SettingData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.base_url = source["base_url"];
	        this.token = source["token"];
	        this.http_proxy = source["http_proxy"];
	    }
	}
	export class UrlsForm {
	    album_id: number;
	    urls: string;
	
	    static createFrom(source: any = {}) {
	        return new UrlsForm(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.album_id = source["album_id"];
	        this.urls = source["urls"];
	    }
	}
	export class UrlsList {
	    page: number;
	    limit: number;
	
	    static createFrom(source: any = {}) {
	        return new UrlsList(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.page = source["page"];
	        this.limit = source["limit"];
	    }
	}

}

