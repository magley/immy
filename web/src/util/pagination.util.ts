import type { ApiResponse } from "@/api/http";
import type { AxiosError, AxiosResponse } from "axios";


export class Paginator<T> {
	perPage = 10;
	page = 1;
	pagesTotal = 0;
	pagesNav: number[] = [];
	loading: boolean = false;

	func: ((offset: number, limit: number) => Promise<AxiosResponse<ApiResponse<T>>>);

	constructor(func: (offset: number, limit: number) => Promise<AxiosResponse<ApiResponse<T>>>, perPage: number = 10) {
		this.perPage = perPage;
		this.func = func;
	}

	async getItems(): Promise<AxiosResponse<ApiResponse<T>>> {
		if (this.page < 1) this.page = 1;
		if (this.page > this.pagesTotal && this.pagesTotal > 0) this.page = this.pagesTotal;
		
		const offset = (this.page - 1) * this.perPage;
		const limit = this.perPage;
		
		this.loading = true;

		return await this.func(offset, limit).then((res) => {
			if (res.data.meta) {
				const meta = res.data.meta;

				console.log(this.pagesTotal);

				this.page = meta.page;
				this.pagesTotal = meta.total_pages;
				this.pagesNav = [
					this.page - 4, this.page - 3, this.page - 2, this.page - 1,
					this.page - 0,
					this.page + 1, this.page + 2, this.page + 3, this.page + 4, this.page + 5,
				];
				this.pagesNav = this.pagesNav.filter((v) => v >= 1 && v <= this.pagesTotal);
			}

			return res;
		}).catch((err: AxiosError) => {
			console.error(err);
			throw err;
		}).finally(() => {
			this.loading = false;
		});
	}
}