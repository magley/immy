import requests
import json


URL = "http://localhost:8080/api/v1/"


def RESP(resp):
	print(resp.ok)
	print(resp.json())


def create_board(name: str, code: str, desc: str = None):
	resp = requests.post(URL + "boards/", json={
		'name': name, 
		'code': code, 
		'description': desc
	})

	RESP(resp)
	
	
def update_board(board_code: str, fields: dict):
	resp = requests.put(URL + "boards/" + board_code, json=fields)

	RESP(resp)	
	

def create_post(board_code: str, text: str, name: str, options: str):
	dto = {
		'board_code': board_code,
		'subject': '',
		'locked': False,
		'sticky': False,
		'post': {
			'name': name,
			'content': text,
			'filename': '/',
			'options': options
		}
	}
	
	print(name)
	
	resp = requests.post(URL + "threads/", json=dto)
	
	RESP(resp)
	
	
	
if __name__ == "__main__":	
	#create_board("Anime & Manga", "a", "Board for discussing anime and manga. Konichiwa, dude!")
	# update_board("2", {"code": "1"})
	# create_board("2", "2", "jjj")
	# update_board("2", {"code": "1"})
	
	create_post('g', 'bla', '', '')
	
	create_post('g', 'bla', 'u', '')
	create_post('g', 'bla', 'u#', '')
	create_post('g', 'bla', 'u##', '')
	create_post('g', 'bla', 'u###', '')
	
	create_post('g', 'bla', 'u#p', '')
	create_post('g', 'bla', 'u##p', '')
	create_post('g', 'bla', 'u###p', '')
	
	create_post('g', 'bla', '#p', '')
	create_post('g', 'bla', '##p', '')
	create_post('g', 'bla', '##$p', '')
	
	create_post('g', 'bla', '#', '')
	create_post('g', 'bla', '##', '')
	create_post('g', 'bla', '##$', '')
