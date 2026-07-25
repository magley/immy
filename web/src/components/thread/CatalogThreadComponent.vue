<script setup lang="ts">
import type { BoardDTO } from '@/api/board.api';
import { CdnAPI } from '@/api/cdn.api';
import type { ThreadForCatalogDTO } from '@/api/thread.api';
import { ThreadToCanonicalForm } from '@/model/thread/thread.model';
import { useRouter } from 'vue-router';


interface Props {
    board: BoardDTO;
    thread: ThreadForCatalogDTO;

    imageSize: number;
    showComment: boolean;

    isShowingHiddenOnly: boolean;
    hiddenThreads: string[];
    pinnedThreadIDs: string[]
}

const props = defineProps<Props>();
const emits = defineEmits(["onClickMenuArrow", "setHidden", "toggleHidden", "setPinned", "togglePinned"]);
const router = useRouter();

const isHidden = () => {
	if (props.board == undefined) return false;
	return props.hiddenThreads.indexOf(ThreadToCanonicalForm(props.board.code, props.thread.post.num)) != -1;
}

const isPinned = (): boolean => {
    if (props.board == undefined) return false;
    return props.pinnedThreadIDs.indexOf(ThreadToCanonicalForm(props.board.code, props.thread.post.num)) != -1;
}

const hideThread = () => {
    emits("setHidden", props.thread, true);
}

const unhideThread = () => {
    emits("setHidden", props.thread, false);
}

const toggleHiddenThread = () => {
    emits("toggleHidden", props.thread);
}

const pinThread = () => {
    emits("setPinned", props.thread, true);
}

const unpinThread = () => {
    emits("setPinned", props.thread, false);
}

const togglePin = () => {
    emits("togglePinned", props.thread);
}

const getDynamicImageStyle = (): any => {
    if (props.thread.post.img_height > props.thread.post.img_width) {
        return { height: (140 * props.imageSize / 100) + 'px', };
    } else {
        return { width: (160 * props.imageSize / 100) + 'px', };
    }
}

const onClickThreadImage = (e: MouseEvent) => {
    if (e.shiftKey) {
        toggleHiddenThread();
    }
    else if (e.altKey) {
        togglePin();
    }
    else {
        if (props.board) {
            const destination = `/${props.board.code}/thread/${props.thread.thread.post_num}`;
            router.push(destination);
        }
    }
}

const onClickMenuArrow = () => {
    emits("onClickMenuArrow", props.thread);
}

</script>

<template>
<span class="catalog-thread-container">
    <span v-if='isShowingHiddenOnly == isHidden()' class="catalog-post">
        <div class="image-container">
            <a href="#" @click.prevent="onClickThreadImage">
                <img v-if="!thread.post.filename && thread.post.md5"
                    :class="{pinned: isPinned()}"
                    :src="CdnAPI.GetPublicURI('file_deleted.png')"
                >
                <!-- Spoiler -->
                <img v-else-if="board.config.allow_spoilers && thread.post.spoiler"
                    :class="{pinned: isPinned()}"
                    class="spoiler"
                    :style="getDynamicImageStyle()"
                    :src="CdnAPI.GetSpoilerURI(board.config.spoiler_image)"
                >
                <!-- Regular thumbnail -->
                <img v-else
                    :class="{pinned: isPinned()}"
                    :style="getDynamicImageStyle()"
                    :src="CdnAPI.GetPostImageThumbnailURI(thread.post)"
                >
            </a>
            <div class="inside-image">
                <img src="/icons/sticky.png" v-if="thread.thread.sticky" title="Sticky"/>
                <img src="/icons/lock.png" v-if="thread.thread.locked" title="Locked"/>
                <a v-if="isPinned()" href="#" @click.prevent="unpinThread()">
                    <img src="/icons/pin.png" title="Pinned - click to unpin" />
                </a>
                <a v-if="isHidden()" href="#" @click.prevent="unhideThread()">
                    <img src="/icons/visible.png" title="Hidden - click to unhide" />
                </a>
            </div>
        </div>

        <br />

        <span class="stats">
            <abbr title="Number of replies">R</abbr>: <strong>{{ thread.stats.post_count }}</strong>
            /
            <abbr title="Number of images">I</abbr>: <strong>{{ thread.stats.image_count }}</strong>
            /
            <abbr title="Number of users">U</abbr>: <strong>{{ thread.stats.user_count }}</strong>
            /
            <a href="#" class="no-underline" :id="`thread-arrow-${thread.thread.post_num}`" @click.prevent="onClickMenuArrow">▶</a>
        </span>
        <br />

        <span class="body" v-if="showComment">
            <template v-if="thread.thread.subject"><span class="subject">{{thread.thread.subject}}</span>: </template>
            <span v-if="thread.post" class="content">
                {{ thread.post.content }}
                <span v-if="thread.post.html">
                    <br/>
                    <span v-html="thread.post.html"></span>
                </span>
            </span>
        </span>
    </span>
</span>
</template>

<style scoped lang="css">
.catalog-post {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    overflow: hidden;
    /*background-color: white;
    border: 1px solid black;*/
    padding: 5px;


    .image-container {
        position: relative;

        img {
        /*	max-width: 160px;      <---  Programatically computed
            max-height: 140px;     <---  Programatically computed */
            object-fit: contain;
            box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2), 0 6px 20px 0 rgba(0, 0, 0, 0.19);

            &.pinned {
                border: 4px dashed var(--background-color-accent);
            }

            &.pinned:hover {
                border: 4px dashed var(--banner-title-color);
            }

            &.spoiler {
                width: 160px !important;
                height: 160px !important;
            }
        }

        .inside-image {
            position: absolute;
            top: 0px;
            left: 0px;
            z-index: 1000;
        }
    }

    .stats {
        cursor: help;
        text-align: center;
        align-self: center;
        font-size: small;
    }

    .body {
        text-align: center;

        .subject {
            font-weight: bold;
        }

        .content {
        }
    }
}
</style>