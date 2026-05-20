import { SplitFilename } from "@/util/file.util";
import { ENV } from "./http";
import type { PostDTO } from "./post.api";

export class CdnAPI {
	static GetURI(fname: string): string {
		return `${ENV.CDN}/files/${fname}`;
	}

	static GetPostImageURI(post: PostDTO): string | undefined {
		if (!post.filename) {
			return undefined;
		}
		return this.GetURI(`${post.filename}`);
	}

	static GetPostImageThumbnailURI(post: PostDTO): string | undefined {
		if (!post.filename) {
			return undefined;
		}
		const [base, _] = SplitFilename(post.filename);
		return this.GetURI(`${base}-thumb.jpg`);
	}
}