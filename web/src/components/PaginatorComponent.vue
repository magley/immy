<script setup lang="ts">
	import type { Paginator } from '@/util/pagination.util';

	interface Props {
		paginator: Paginator<any>,
		emptyMessage: string,
	}
	const props = withDefaults(defineProps<Props>(), { emptyMessage:"Nothing to Naviage" });
	const emits = defineEmits(["gotoPage"]);

	const gotoPage = (page: number) => {
		emits("gotoPage", page);
	}
</script>

<template>
	<template v-if="paginator.pagesTotal==0">
		[{{ emptyMessage }}]
	</template>
	<template v-else>
		[<a href="#" @click.prevent="gotoPage(1)">First</a>]&thinsp;
			
		[<template v-if="paginator.page > 1"><a href="#" @click.prevent="gotoPage(paginator.page - 1)">Prev</a></template>
		<template v-else>Prev</template>]&thinsp;

		<span v-for="p, i of paginator.pagesNav">
			<template v-if="p == paginator.page">
				<span>{{ p }} </span>
			</template>
			<template v-else>
				<a href="#" @click.prevent="gotoPage(p)">{{ p }} </a>
			</template>
			<template v-if="i < paginator.pagesNav.length - 1">,</template>&thinsp;
		</span>

		[<template v-if="paginator.page < paginator.pagesTotal"><a href="#" @click.prevent="gotoPage(paginator.page + 1)">Next</a></template>
		<template v-else>Next</template>]&thinsp;

		[<a href="#" @click.prevent="gotoPage(paginator.pagesTotal)">Last</a>]&thinsp;
	</template>
</template>