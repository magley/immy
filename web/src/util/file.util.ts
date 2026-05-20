export const GetFileFromEvent = (e: Event, max_size_bytes: number): [File | null, string | null] => {
	const input = e.target as HTMLInputElement;
	const filesAsArray = Array.from(input?.files || []);

	if (filesAsArray.length == 0) {
		return [null, null];
	}

	const file: File | undefined = filesAsArray[0];
	if (file == undefined) {
		return [null, "No file attached"];
	}

	if (file.size > max_size_bytes) {
		const maxsize_mb = Math.round(max_size_bytes / (1024 * 1024));
		return [null, `File cannot be larger than ${maxsize_mb}MB`];
	}

	return [file, null];
}

export const FileToBase64 = async (file: File): Promise<string> => {
	const arrayBuffer = await file.arrayBuffer();
	const uint8Array = new Uint8Array(arrayBuffer);
	return btoa(String.fromCharCode(...uint8Array));
}

export const SplitFilename = (fname: string): [string, string] => {
	const parts: string[] = fname.split(".");
	return [parts.slice(0, -1).join("."), parts[parts.slice.length - 1]!];
}