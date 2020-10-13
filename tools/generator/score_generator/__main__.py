import os
import time
import json
import random
import requests
import threading
from faker import Faker
from pprint import pprint
from multiprocessing import Pool

BASE_URL = 'http://leaderboard-v2-lb-ecs-tg-584908050.eu-central-1.elb.amazonaws.com'

if 'BASE_URL' in os.environ.keys():
    BASE_URL = os.environ['BASE_URL']

DURATIONS = {}
DURATIONS_LOCK = threading.Lock()
FAKER = Faker()
UNTIL_USERS_REACHED = 1000


def update_statistics(key, r):
    DURATIONS_LOCK.acquire()
    if key not in DURATIONS.keys():
        DURATIONS[key] = {
            'sum': 0,
            'n_requests': 0,
            'avg': 0,
            'last': 0
        }

    stats_dict = DURATIONS[key]
    stats_dict['sum'] += r.elapsed.microseconds
    stats_dict['n_requests'] += 1
    stats_dict['avg'] = stats_dict['sum'] // stats_dict['n_requests']
    stats_dict['last'] = r.elapsed.microseconds
    DURATIONS[key] = stats_dict
    DURATIONS_LOCK.release()


def get_leaderboard():
    r = requests.get('%s/leaderboard' % BASE_URL, params={'page_size': 10})
    update_statistics('leaderboard', r)


def create_user():
    user = {
        'display_name': FAKER.user_name(),
        'country': FAKER.country_code()
    }

    r = requests.post('%s/user/create' % BASE_URL, json=user)
    try:
        assert r.status_code == 201
    except AssertionError:
        print(r.text)
        exit(1)
    update_statistics('create_user', r)

    r = r.json()
    return r['user_id']


def play_game(user_id):
    r = requests.post('%s/score/submit' % BASE_URL, json={
        'score': (random.random() * 100000) + 0.1,
        'user_id': user_id,
        'timestamp': int(time.time())
    })
    try:
        assert r.status_code == 201
    except AssertionError:
        print(r.text)
        exit(1)
    update_statistics('play_game', r)


def generate_data(id):
    to_process = UNTIL_USERS_REACHED // os.cpu_count()
    print('#%d started for %d users.' % (id, to_process))

    USER_IDS = []
    for current_user_index in range(to_process):
        user_id = create_user()
        USER_IDS.append(user_id)

        completion = (current_user_index * 100 // to_process)
        if completion != 0 and completion % 10 == 0:
            get_leaderboard()
            print('#%d: %.2f : %s' %
                  (id, completion, json.dumps(DURATIONS)))

    for current_game_index in range(to_process * 100):
        play_game(random.choice(USER_IDS))
        completion = (current_game_index * 100 // to_process)
        if completion != 0 and completion % 10 == 0:
            get_leaderboard()
            print('#%d: %.2f : %s' %
                  (id, completion, json.dumps(DURATIONS)))

    pprint(DURATIONS)
    return 1


with Pool(os.cpu_count()) as pool:
    print(sum(pool.map(generate_data, [x for x in range(os.cpu_count())])))
