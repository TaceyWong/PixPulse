export namespace converter {
	
	export class ConvertRequest {
	    inputPath: string;
	    outputPath: string;
	    mode: string;
	
	    static createFrom(source: any = {}) {
	        return new ConvertRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.inputPath = source["inputPath"];
	        this.outputPath = source["outputPath"];
	        this.mode = source["mode"];
	    }
	}
	export class ConvertResponse {
	    success: boolean;
	    error: string;
	
	    static createFrom(source: any = {}) {
	        return new ConvertResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.error = source["error"];
	    }
	}

}

