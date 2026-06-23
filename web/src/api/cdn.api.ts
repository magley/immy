import { SplitFilename } from "@/util/file.util";
import { axiosInstanceCDN, ENV } from "./http";
import type { PostDTO } from "./post.api";
import type { AxiosResponse } from "axios";

export class CdnAPI {
	static GetFilesURI(fname: string): string {
		return `${ENV.CDN}/files/${fname}`;
	}

	static GetBannersURI(fname: string): string {
		return this.GetPublicURI(`/banners/${fname}`);
	}

	static GetSpoilerURI(fname: string): string {
		return this.GetPublicURI(`/spoiler/${fname}`);
	}

	static GetPublicURI(fname: string): string {
		return `${ENV.CDN}/public/${fname}`;
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
		return this.GetFilesIn("public/banners/title");
	}

	static async GetBoardBanners(): Promise<string[]> {
		return this.GetFilesIn("public/banners/board");
	}

	static GetTitleBanner = (fname: string[]): string | undefined => {
		if (fname.length == 0) {
			return undefined;
		}
		const i = Math.floor(Math.random() * fname.length);
		return this.GetBannersURI("title/" + fname[i]!);
	}

	static GetBoardBanner = (fname: string[]): string | undefined => {
		if (fname.length == 0) {
			return undefined;
		}
		const i = Math.floor(Math.random() * fname.length);
		return this.GetBannersURI("board/" + fname[i]!);
	}

	static GetFilesIn = (folder: string) => {
		return axiosInstanceCDN.get(folder).then((res) => {
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
}