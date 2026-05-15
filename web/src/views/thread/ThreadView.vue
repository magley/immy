<script setup lang="ts">
	import { ref, onMounted } from 'vue';
	import { useRoute, useRouter } from "vue-router";
	import { BoardAPI, BoardDTO, CreateBoardDTO, UpdateBoardDTO } from "@/api/board.api.ts";
	import { ThreadAPI, ThreadDTO, CreateThreadDTO, UpdateThreadDTO } from "@/api/thread.api.ts";
	import { PostAPI, PostDTO, CreatePostForThreadDTO, CreatePostDTO, UpdatePostDTO } from "@/api/post.api.ts";
	
	const board = ref<BoardDTO | null>(null);

	const route = useRoute();
	const router = useRouter();
	
	const thread = ref<ThreadDTO>(null);
	const posts = ref<PostDTO[]>([]);
	
	const replyDTO = ref<CreatePostDTO>({});
	const replyError = ref<string | undefined>(undefined);
	
	onMounted(() => {
		const board_code: string = route.params.board_code;
		loadBoard(board_code);
	});
	
	const onSubmitReply = () => {
		console.log(replyDTO.value.content);
		replyError.value = undefined;
		replyDTO.value.thread_id = thread.value.id;

		PostAPI.CreatePost(replyDTO.value).then((res: AxiosResponse<ApiResponse<PostDTO>>) => {
			loadThread(board.value.code, thread.value.post_num);
		}).catch((err: AxiosError) => {
			replyError.value = "Could not post reply";
			console.error(err);
		});
	}
	
	const loadBoard = (boardCode: string) => {
		BoardAPI.GetBoard(boardCode).then((res: AxiosResponse<ApiResponse<BoardDTO>>) => {
			board.value = res.data.data;
			
			const thread_num: number = route.params.thread_num;
			loadThread(boardCode, thread_num);
		}).catch((err: AxiosError) => {
			router.push("/");
		});
	}
	
	const loadThread = (board_code: string, thread_num: number) => {
		ThreadAPI.GetFullThreadByNum(board_code, thread_num).then((res: AxiosResponse<ApiResponse<ThreadFullDTO>>) => {
			const dto: ThreadFullDTO = res.data.data;
			thread.value = dto.thread;
			posts.value = dto.posts;
			
			for (let p of posts.value) {
				console.log(p.content);
			}
		}).catch((err: AxiosError) => {
			console.error(err);
		});
	}
	
	const getPostTimeReadable = (dateStr: string) => {
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
	
	const onClickPostNumber = (num: number) => {
		if (replyDTO.value.content == undefined) {
			replyDTO.value.content = "";
		}
		replyDTO.value.content += `>>${num}\n`;
	}
	
	const processedPostContent = (text: string) => {
		// replace >>[num] with >><a href=#p[num]>>>[num]</a>
		text = text.replace(/>>(\w+)/g, '<a href="#p$1" class="postRef">&gt;&gt;$1</a>');
		return text;
	}
</script>

<template>
	<template v-if="board">
		Thread {{ route.params.thread_id }}
		/{{board.code}}/ - {{board.name}}
		<br/>
		
		<RouterLink :to="`/${route.params.board_code}`">[Return]</RouterLink>
		<RouterLink :to="`/${route.params.board_code}/catalog`">[Catalog]</RouterLink>
		<br/>
		
		
		<!-- New reply -->
		<form @submit.prevent="onSubmitReply">
			<input type=text placeholder="Subject" v-model="replyDTO.subject"/><br/>
			<input type=text placeholder="Name" v-model="replyDTO.name"/><br/>
			<input type=text placeholder="Options" v-model="replyDTO.options"/><br/>
			<textarea placeholder="Text..." v-model="replyDTO.content"/><br/>
			Files not implemented yet...
			<br/>
			<button type=submit>Post reply</button>
			
			<template v-if="replyError">
				<div/>
				<span class="error">{{replyError}}</span>
			</template>
		</form>
		
		<template v-if="thread">
			<div :id="`p${post.num}`" class="postContainer" v-for="post, i of posts">
				<span class="sideArrows"> &gt;&gt; </span>
				<span class="post">
					<div class="post-header">
						<!-- >> Anonymous 05/15/26(Fri)16:12:37 No.221512455 ▶ >>221512482  -->
						<span class="subject" v-if="thread.subject && thread.post_num == post.num">{{ thread.subject }}</span>
						<span class="username">{{ post.name ? post.name : "Anonymous" }}</span>
						<span class="date">{{ getPostTimeReadable(post.created_at) }}</span>
						<span class="postno">No.</span>
						<span class="postnum"><a @click.prevent="onClickPostNumber(post.num)" href="#" class="postNumLink">{{ post.num }}</a></span>
						<span class="dropdown">&#9654;</span>
					</div>
					
					<div class="post-body" v-html="processedPostContent(post.content)"></div>
				</span>
			</div>
		</template>
	</template>
</template>

<style scoped>	
	html {
		background-color: #EEF2FF;
	}
	
	.postContainer {
		display: flex;
		gap: 10px;
		margin-top: 0.2em;

		.sideArrows {
			
		}
		
		.post {
			flex-grow: 1;
			
			background-color: #D6DAF0;
			padding-top: 0.25em;
			padding-bottom: 1em;
			padding-left: 1em;
			padding-right: 1em;
			
			.opPost {
				
			}
			
			.post-header {
				margin-bottom: 1em;
				
				span {
					margin-right: 0.25em;
				}
			}
			
			.post-body {
				margin-left: 1em;
				white-space: pre-wrap; 
  				word-wrap: break-word;
			}		
		}
	}
	
	.username {
		color: #157743;
		font-weight: bolder;
	}
	
	.subject {
		color: #0F0C5D;
		font-weight: bolder;
	}
	
	.postNumLink {
		color: black;
		text-decoration: none;
	}
	
	.postNumLink:hover {
		color: #DD0000;
	}
</style>