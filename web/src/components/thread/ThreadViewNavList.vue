<script setup lang="ts">
	import { type ThreadStats } from "@/model/thread/thread.model.ts";

	export interface ThreadNavProps {
		board_code: string,
		jump_to_id: string | null,
		jump_to_label: string | null,
		thread_stats: ThreadStats,
		sticky: boolean,
		locked: boolean,
		showCenterElements: boolean,
		//
		autoTimer: number,
		isAutoTimerUsed: boolean,
	}

	const props = defineProps<ThreadNavProps>();
	const emit = defineEmits(['autoTimerToggled', 'updateClicked', 'openedReplyBox', 'openedGalleryMode']);
	
	const onSetAutoTimer = (ev: Event) => {
		emit('autoTimerToggled', (ev.target as HTMLInputElement).checked);
	}

	const doUpdate = () => {
		emit('updateClicked');
	}

	const openedReplyBox = () => {
		emit('openedReplyBox');
	}

	const openedGalleryMode = () => {
		emit('openedGalleryMode');
	}
</script>

<template>
<hr />
<nav>
	<span class="left">
		[<RouterLink class="link" :to="`/${props.board_code}`">Return</RouterLink>]
		[<RouterLink class="link" :to="`/${props.board_code}/catalog`">Catalog</RouterLink>]
		<template v-if="props.jump_to_id && props.jump_to_label">
			[<a class="link" :href="`#${props.jump_to_id}`">{{ props.jump_to_label }}</a>]
		</template>
		[<a class="link" href="#" @click.prevent="doUpdate">Update</a>]
		[<input type="checkbox" :checked="props.isAutoTimerUsed" @change="onSetAutoTimer" name="auto"><label for="auto"> Auto</label>]
		<template v-if="props.isAutoTimerUsed">{{ props.autoTimer }}</template>
		|
		[<a class="link" href="#" @click.prevent="openedGalleryMode">Gallery</a>]
	</span>
	<span class="middle" v-if="showCenterElements">
		[<a class="link" href="#" @click.prevent="openedReplyBox">Post a Reply</a>]
	</span>
	<span class ="right">
		<template v-if="sticky">Sticky / </template>
		<template v-if="locked">Locked / </template>
		<span class="tooltip">{{ props.thread_stats.posts }}
			<span class="tooltiptext">Replies</span>
		</span> / 
		<span class="tooltip">{{ props.thread_stats.images }}
			<span class="tooltiptext">Images</span>
		</span> / 
		<span class="tooltip">{{ props.thread_stats.posters }}
			<span class="tooltiptext">Users</span>
		</span> / 
		<span class="tooltip">{{ props.thread_stats.page }}
			<span class="tooltiptext">Page</span>
		</span>
	</span>
</nav>
<hr />
</template>

<style scoped>
	nav {
		display: flex;
		justify-content: space-between;

		.left {
			flex: 1;
		}

		.middle {
			flex: 1;
			text-align: center;
		}

		.right {
			flex: 1;
			cursor: default;
			text-align: right;
		}
	}
</style>