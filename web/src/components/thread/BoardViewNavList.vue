<script setup lang="ts">
	import { ThreadSortModeInCatalog } from "@/model/thread/thread.model.ts";

	interface BoardNavProps {
		board_code: string;
		jump_to_id: string;
		jump_to_label: string;

		sortBy: ThreadSortModeInCatalog;
		imageSize: number;
		showComment: boolean;
	}

	const props = defineProps<BoardNavProps>();
	const emit = defineEmits(['sortChanged', 'imageSizeChanged', 'showCommentChanged']);

	const onSortFieldChange = (ev: Event) => {
		emit('sortChanged', (ev.target as HTMLSelectElement).value as ThreadSortModeInCatalog);
	}
	const onImageSizeChange = (ev: Event) => {
		emit('imageSizeChanged', Number((ev.target as HTMLInputElement).value));
	}
	const onShowCommentChange = (ev: Event) => {
		emit('showCommentChanged', (ev.target as HTMLInputElement).checked);
	}
</script>

<template>
<nav>
	<span class="left">
		[<RouterLink class="link" :to="`/${props.board_code}`">Return</RouterLink>]
		[<RouterLink class="link" :to="`/${props.board_code}/archive`">Archive</RouterLink>]
		<template v-if="props.jump_to_id && props.jump_to_label">
			[<a class="link" :href="`#${props.jump_to_id}`">{{ props.jump_to_label }}</a>]
		</template>
		[<RouterLink class="link" :to="`/${props.board_code}/catalog`">Refresh</RouterLink>]
	</span>
	<span class ="right">
		<span class="setting">
			<label for="sort-field">Sort by: </label>
			<select id="sort-field" :value="props.sortBy" @change="onSortFieldChange">
			  <option :value="ThreadSortModeInCatalog.BumpOrder">Bump Order</option>
			  <option :value="ThreadSortModeInCatalog.LastReply">Last Reply</option>
			  <option :value="ThreadSortModeInCatalog.CreationDate">Creation Date</option>
			  <option :value="ThreadSortModeInCatalog.ReplyCount">Reply Count</option>
			  <option :value="ThreadSortModeInCatalog.ImageCount">Image Count</option>
			  <option :value="ThreadSortModeInCatalog.UserCount">User Count</option>
			</select>
		</span>
		<span class="setting">
	  		<label for="image-size">Image size: </label>
			<input type="range" id="image-size" min="0" max="200" :value="props.imageSize" @input="onImageSizeChange" />
		</span>
		<span class="setting">
			<label for="show-comment">Show comment: </label>
			<input type="checkbox" id="show-comment" :checked="props.showComment" @change="onShowCommentChange" />
		</span>
	</span>
</nav>
</template>

<style scoped>
	nav {
		display: flex;
		justify-content: space-between;

		.right {
			cursor: default;
		}
	}

	.right {
		.setting {
			margin-left: 1em;
		}
	}
</style>