import { SplitFilename } from "@/util/file.util";
import { axiosInstanceCDN, ENV } from "./http";
import type { PostDTO } from "./post.api";
import type { AxiosResponse } from "axios";

export class CdnAPI {
	static GetFilesURI(fname: string): string {
		return `${ENV.CDN}/files/${fname}`;
	}

	static GetBannersURI(fname: string): string {
		return `${ENV.CDN}/banners/${fname}`;
	}

	static GetPostImageURI(post: PostDTO): string | undefined {
		if (!post.filename) {
			return undefined;
		}
		return this.GetFilesURI(`${post.filename}`);
	}

	static GetPostImageThumbnailURI(post: PostDTO): string | undefined {
		if (!post.filename) {
			return undefined;
		}
		const [base, _] = SplitFilename(post.filename);
		return this.GetFilesURI(`${base}-thumb.jpg`);
	}

	static async GetTitleBanners(): Promise<string[]> {
		return axiosInstanceCDN.get("banners/title").then((res) => {
			/*
			[
				{ "name":"b_01.png", "type":"file", "mtime":"Fri, 12 Jun 2026 12:40:59 GMT", "size":11688 },
				{ "name":"g_01.gif", "type":"file", "mtime":"Fri, 12 Jun 2026 12:44:04 GMT", "size":12846 }
			]
			*/
			return (res.data! as unknown as any[]).map((r: any) => r.name!);
		}).catch((err) => {
			throw err;
		});
	}

	static async GetBoardBanners(): Promise<string[]> {
		return axiosInstanceCDN.get("banners/board").then((res) => {
			return (res.data! as unknown as any[]).map((r: any) => r.name!);
		}).catch((err) => {
			throw err;
		});
	}

	static GetTitleBanner = (fname: string[]): string | undefined => {
		if (fname.length == 0) {
			return undefined;
		}
		const i = Math.floor(Math.random() * fname.length);
		console.log(fname, i);
		return this.GetBannersURI("title/" + fname[i]!);
	}

	static GetBoardBanner = (fname: string[]): string | undefined => {
		if (fname.length == 0) {
			return undefined;
		}
		const i = Math.floor(Math.random() * fname.length);
		console.log(fname, i);
		return this.GetBannersURI("board/" + fname[i]!);
	}
}