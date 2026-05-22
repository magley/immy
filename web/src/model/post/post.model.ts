export const GetPostTimeReadable = (dateStr: string) => {
	var date: Date = new Date(dateStr);

	const getDayOfWeek = (date: Date) => {
		switch (date.getDay()) {
		case 1: return "Mon";
		case 2: return "Tue";
		case 3: return "Wed";
		case 4: return "Thu";
		case 5: return "Fri";
		case 6: return "Sat";
		case 7: return "Sun";
		default: return "???";
		}
	}

	const getDateStr = (date: Date) => {
		const d = date.getDate();
		const m = date.getMonth() + 1;
		const y = date.getFullYear();

		const dd = String(d).padStart(2, '0');
		const mm = String(m).padStart(2, '0');
		const yy = String(y).padStart(2, '0').substring(2);

		return `${dd}/${mm}/${yy}`;
	}

	const getTimeStr = (date: Date) => {
		const h = date.getHours();
		const m = date.getMinutes();
		const s = date.getSeconds();

		const hh = String(h).padStart(2, '0');
		const mm = String(m).padStart(2, '0');
		const ss = String(s).padStart(2, '0');
		return `${hh}:${mm}:${ss}`;
	}

	return `${getDateStr(date)} (${getDayOfWeek(date)})${getTimeStr(date)}`;
}


export type PostTextToken = {
	kind: "text";
	text: string;
};

export type PostLinkToken = {
	kind: "link";
	text: string;
	local: boolean,
	fail: boolean
};

export type PostToken = PostTextToken | PostLinkToken;

export const ParsePostTokens = (text: string): PostToken[] => {
	return text.split(/(\s+|##\w+|\S+)/g).map(word => {
		if (word.startsWith(">>")) {
			return {
				kind: "link",
				text: word,
			} as unknown as PostTextToken;
		}
		return {
			kind: "text",
			text: word,
			local: true,
			fail: true,
		} as unknown as PostLinkToken;
	});
}

export interface PostImageData {
	postId: number,
	expanded: boolean,
	width: number,
	height: number,
}