import type { UserRole } from "@/api/user.api";

export const AddRangeNoDuplicates = <T>(target: T[], range: T[]): T[] => {
	return range.reduce((acc, item) => {return acc.includes(item) ? acc : [...acc, item] }, [...target]);
}

export const UniqueArray = <T>(arr: T[]): T[] => {
	return arr.filter((x, i, a) => a.indexOf(x) == i);
}


export const GetElipsisString = (input: string, lenIncludingDots: number = 10 + 3) => {
	if (!input) return input;

	if (input.length > lenIncludingDots) {
		return `${input.substring(0, lenIncludingDots)}...`;
	} else {
		return input;
	}
}

export const IsAlphaNumeric = (str: string) =>  {
	var code, i, len;

	for (i = 0, len = str.length; i < len; i++) {
		code = str.charCodeAt(i);
		if (!(code > 47 && code < 58) && !(code > 64 && code < 91) && !(code > 96 && code < 123)) {
			return false;
		}
	}
	return true;
}

export const StripSlashes = (str: string) => {
	if (str.startsWith("/") && str.endsWith("/")) {
		return str.slice(1, -1);
	}
	return str;
}

export const GetPublicIdColorBackground = (user_id: string): string => {
	var hash = 0;
	if (user_id.length === 0) {
		return "#FFFFFF";
	}
	for (var i = 0; i < user_id.length; i++) {
		hash = user_id.charCodeAt(i) + ((hash << 5) - hash);
		hash = hash & hash;
	}

	var color = '#';
	for (var i = 0; i < 3; i++) {
		var value = (hash >> (i * 8)) & 255;
		color += ('00' + value.toString(16)).substr(-2);
	}

	return color;
}

export const GetPublicIdColorForeground = (user_id: string): string => {
	const bgIsLIght: boolean = IsHexColorLight(GetPublicIdColorBackground(user_id));

	if (bgIsLIght) {
		return "#000000";
	} else {
		return "#FFFFFF";
	}
}

export const IsHexColorLight = (color: string): boolean => {
	if (color.length == 6) {
		return IsHexColorLight("#" + color);
	}

	if (color.length == 7) {
		const rgb: number[] = [
			parseInt(color.substring(1, 3), 16),
			parseInt(color.substring(3, 5), 16),
			parseInt(color.substring(5), 16),
		];
		const luminance =
		(0.2126 * rgb[0]!) / 255 +
		(0.7152 * rgb[1]!) / 255 +
		(0.0722 * rgb[2]!) / 255;
		return luminance > 0.5;
	}
	return false
}

export const DateFromDuration = (input: string): Date => {
	let value = parseInt(input, 10);
	if (isNaN(value)) value = 0;

	const unit = input
	.slice(String(value).length)
	.trim()
	.toLowerCase();

	const SECOND = 1000;
	const MINUTE = 60 * SECOND;
	const HOUR = 60 * MINUTE;
	const DAY = 24 * HOUR;
	const WEEK = 7 * DAY;
	const YEAR = 365 * DAY;

	const unitMs: Record<string, number> = {
		min: MINUTE,
		minute: MINUTE,
		minutes: MINUTE,
		h: HOUR,
		hour: HOUR,
		hours: HOUR,
		d: DAY,
		day: DAY,
		days: DAY,
		w: WEEK,
		week: WEEK,
		weeks: WEEK,
		y: YEAR,
		year: YEAR,
		years: YEAR,
	};

	const unitVal = unitMs[unit] ?? 0;
	let plusDelta = value * unitVal;
	if (plusDelta < 0) {
		plusDelta = 0;
	}

	return new Date(Date.now() + plusDelta);
}

export const GetTimeDifferenceBasic = (from: Date, to: Date): string => {
	const diffMs = Math.abs(to.getTime() - from.getTime());

	const second = 1000;
	const minute = 60 * second;
	const hour = 60 * minute;
	const day = 24 * hour;
	const month = 30 * day;
	const year = 365 * day;

	if (diffMs >= year) {
		return `${Math.round(diffMs / year)} year(s)`;
	}

	if (diffMs >= month) {
		return `${Math.round(diffMs / month)} month(s)`;
	}

	if (diffMs >= day) {
		return `${Math.round(diffMs / day)} day(s)`;
	}

	if (diffMs >= hour) {
		return `${Math.round(diffMs / hour)} hour(s)`;
	}

	if (diffMs >= minute) {
		return `${Math.round(diffMs / minute)} minute(s)`;
	}

	return `${Math.round(diffMs / second)} second(s)`;
}

export const RemoveLoginCredentials = () => {
	localStorage.removeItem("username");
	localStorage.removeItem("id");
	localStorage.removeItem("role");
	localStorage.removeItem("jwt");
}

// Returns `false` if `jwt` is `null`.
export const IsJwtExpired = (jwt: string | null): boolean => {
	if (jwt == null) return false;

	const parts: string[] = jwt.split(".");
	const bodyStr: string = atob(parts[1]!);

	interface JSONBody {
		username: string;
		id: number;
		role: UserRole;
		iat: number;
		exp: number;
	}

	const body: JSONBody = JSON.parse(bodyStr);
	const expire_ms = body.exp * 1000;
	return expire_ms < new Date().getTime();
}