import { ENV } from "./http";

export class CdnAPI {
	static GetURI(fname: string): string {
		return `${ENV.CDN}/files/${fname}`;
	}
}