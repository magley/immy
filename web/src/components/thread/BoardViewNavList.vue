<script setup lang="ts">
	import { ref, onMounted, defineProps, defineEmits } from 'vue';
	import { type ThreadStats, type ThreadNavProps } from "@/model/thread/thread.model.ts";

	interface BoardNavProps {
		board_code: string;
		jump_to_id: string;
		jump_to_label: string;
	}

	const props = defineProps<BoardNavProps>();
	const emit = defineEmits(['sortChanged', 'imageSizeChanged', 'showCommentChanged']);

	const sortBy = ref<string>("creationDate");
	const imageSize = ref<number>(100);
	const showComment = ref<boolean>(true);

	const onSortFieldChange = () => {
		emit('sortChanged', sortBy.value);
	}
	const onImageSizeChange = () => {
		emit('imageSizeChanged', imageSize.value);
	}
	const onShowCommentChange = () => {
		emit('showCommentChanged', showComment.value);
	}

</script>

<template>
<nav>
	<span class="left">
		[<RouterLink class="link" :to="`/${props.board_code}`">Return</RouterLink>]
		<template v-if="props.jump_to_id && props.jump_to_label">
			[<a class="link" :href="`#${props.jump_to_id}`">{{ props.jump_to_label }}</a>]
		</template>
		[<RouterLink class="link" :to="`/${props.board_code}/catalog`">Refresh</RouterLink>]
	</span>
	<span class ="right">
		<span class="setting">
			<label for="sort-field">Sort by: </label>
			<select id="sort-field" v-model="sortBy" @change="onSortFieldChange">
			  <option value="bumpOrder">Bump Order</option>
			  <option value="lastReply">Last Reply</option>
			  <option value="creationDate">Creation Date</option>
			  <option value="replyCount">Reply Count</option>
			  <option value="imageCount">Image Count</option>
			  <option value="userCount">User Count</option>
			</select>
		</span>
		<span class="setting">
	  		<label for="image-size">Image size: </label>
			<input type="range" id="image-size" min="0" max="200" v-model="imageSize" @input="onImageSizeChange" />
		</span>
		<span class="setting">
			<label for="show-comment">Show comment: </label>
			<input type="checkbox" id="show-comment" v-model="showComment" @change="onShowCommentChange" />
		</span>
	</span>
</nav>
</template>

<style scoped>
	RouterLink, a {
		color: #34345c;
	}

	RouterLink:hover, a:hover {
		color: #DD0000;
	}

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