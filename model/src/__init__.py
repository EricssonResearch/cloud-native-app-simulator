"""
Copyright 2021 Ericsson AB

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

from flask import Flask
from src.controller import routes
from src.service import logger

def create_app(script_info=None):
    app = Flask(__name__)
    app.config.from_object("src.config.DevelopmentConfig")

    app.register_blueprint(routes.simple_page)
    app.shell_context_processor({"app" : app})
    app.run(threaded=True, processes=4)

    return app